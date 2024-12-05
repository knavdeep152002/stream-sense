package room

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/knavdeep152002/stream-sense/internal/constants"
	"github.com/knavdeep152002/stream-sense/internal/models"
	"github.com/knavdeep152002/stream-sense/internal/utils"
	"github.com/knavdeep152002/stream-sense/internal/utils/concurrency"
	"gorm.io/gorm"
)

type RoomHandler struct {
	DB                    *gorm.DB
	roomClientHandlerPool *concurrency.Pool[*RoomClientHandler]
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	rh := &RoomHandler{
		DB:                    db,
		roomClientHandlerPool: concurrency.NewPool[*RoomClientHandler](),
	}
	return rh
}

// @Summary		Create a room
// @Description	Create a room
// @Tags			room
// @Accept			json
// @Produce		json
// @Param			videoId	path		string	true	"Video ID"
// @Success		200		{object}	roomResult
// @Security		Bearer
// @Router			/room/{videoId} [post]
func (rh *RoomHandler) CreateRoom(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	roomModel := &models.Room{
		VideoID: c.Param("videoId"),
		Active:  true,
		AdminID: userID,
	}
	craetedRoom, err := rh.createRoomInDB(roomModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomResult := &roomResult{
		RoomId:   craetedRoom.ID,
		VideoID:  craetedRoom.VideoID,
		RoomLink: fmt.Sprintf("%d@%s", craetedRoom.ID, craetedRoom.VideoID),
	}
	// register room in connection pool buffer
	c.JSON(http.StatusOK, gin.H{"data": roomResult})
}

// @Summary		Get active rooms
// @Description	Get active rooms
// @Tags			room
// @Accept			json
// @Produce		json
// @Success		200	{object}	roomResult
// @Security		Bearer
// @Router			/rooms [get]
func (rh *RoomHandler) GetActiveRooms(c *gin.Context) {
	var rooms []models.Room
	userID := c.GetUint("userID")
	// get current user's rooms in which he is admin and room is active
	err := rh.DB.Where("admin_id=? AND active=?", userID, true).Find(&rooms).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomResultList := make([]roomResult, 0)
	for _, room := range rooms {
		roomResultList = append(roomResultList, roomResult{
			RoomId:   room.ID,
			VideoID:  room.VideoID,
			RoomLink: fmt.Sprintf("%d@%s", room.ID, room.VideoID),
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": roomResultList})
}

// @Summary		Serve video
// @Description	Serve video
// @Tags			room
// @Accept			json
// @Produce		json
// @Param			videoId	path	string	true	"Video ID"
// @Success		200		{file}	file	"ok"
// @Security		Bearer
// @Router			/video/{roomLink} [get]
func (rh *RoomHandler) ServeVideo(c *gin.Context) {
	roomLink := c.Param("roomLink")
	log.Println("Hello")
	if roomLink == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}
	log.Println("Room Link: ", roomLink)
	roomLink = strings.Split(strings.Trim(roomLink, "/"), "/")[0]
	log.Println("Room Link: ", roomLink)
	room, err := rh.getRoomAndVideoID(strings.Trim(roomLink, "/"))
	if err != nil {
		log.Println("Error in getting room and video id: ", err)
		return
	}
	rch, err := rh.getOrCreateRoomClient(roomLink, room.VideoID, room.AdminID)
	if err != nil {
		return
	}
	hlsDir := path.Join(utils.SegmentDir, rch.VideoId)
	if _, err := os.Stat(hlsDir); err == nil {
		// Strip the prefix part and use the remaining path to serve files
		prefix := fmt.Sprintf("%s/api/v1/serve/%s", constants.PATH_PREFIX, roomLink)
		http.StripPrefix(prefix, http.FileServer(http.Dir(hlsDir))).ServeHTTP(c.Writer, c.Request)
	} else {
		c.String(http.StatusNotFound, "File or directory not found")
	}
}

func (rh *RoomHandler) JoinRoom(c *gin.Context) {
	log.Println("Joining room")
	roomLink := c.Param("roomLink")
	userID := c.GetUint("userID")
	if roomLink == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}
	// get the room id and video id
	room, err := rh.getRoomAndVideoID(roomLink)
	if err != nil {
		log.Println("Error in getting room and video id: ", err)
		return
	}
	err = rh.addUserToRoom(room.ID, userID)
	if err != nil {
		log.Println("Error in adding user to room: ", err)
		return
	}
	// upgrade the connection to websocket
	conn, _, _, _ := ws.UpgradeHTTP(c.Request, c.Writer)
	// defer conn.Close()
	outChan := make(chan *RoomStreamResponse, 10)
	inChan := make(chan *RoomController, 10)
	rch, err := rh.getOrCreateRoomClient(roomLink, room.VideoID, room.AdminID)
	if err != nil {
		conn.Close()
		go rh.removeUserFromRoom(room.ID, userID)
		return
	}
	outChanId, err := rch.Publisher.RegisterObserver(outChan)
	var wg sync.WaitGroup
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	wg.Add(2)
	rch.Observer.Observe(inChan)
	rch.OutChan <- &RoomStreamResponse{
		RoomID:    room.ID,
		Message:   fmt.Sprintf("User %d joined the room", userID),
		MsgType:   "Info",
		Status:    rch.state,
		TimeStamp: rch.TimeStamp,
	}
	go func() {
		defer func() {
			log.Println("Connection Closed")
			rh.roomClientHandlerPool.Unset(roomLink)
			conn.Close()
			wg.Done()
		}()
		for rch.Ready {
			msg, err := wsutil.ReadClientText(conn)
			if err != nil {
				log.Println("Error in reading from websocket: ", err)
				return
			}
			log.Println("Message: ", string(msg))
			var control RoomController
			err = json.Unmarshal([]byte(msg), &control)
			if err != nil {
				log.Println("Error in unmarshalling control: ", err)
				return
			}
			control.UserID = userID
			control.RoomID = room.ID
			if control.UserID != room.AdminID && control.MsgType == "control" {
				rch.OutChan <- &RoomStreamResponse{
					RoomID:  room.ID,
					MsgType: "Error",
					UserID:  userID,
					Message: fmt.Sprintf("User %d is not authorized to control the room", userID),
				}
			} else {
				if control.MsgType == "control" || control.MsgType == "current_status" {
					inChan <- &control
				}
			}

		}
	}()

	go func() {
		defer func() {
			log.Println("Closing output channel")
			rch.Publisher.DeregisterObserver(outChanId)
			conn.Close()
			wg.Done()
		}()
		for rch.Ready {
			select {
			case out := <-outChan:
				byteResp, _ := json.Marshal(out)
				err := wsutil.WriteServerText(conn, byteResp)
				if err != nil {
					log.Println("Error in writing to websocket: ", err)
					return
				}
			case <-time.After(time.Second):
			}

		}
	}()
	wg.Wait()
}
