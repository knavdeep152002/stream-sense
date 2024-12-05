package room

import (
	"errors"
	"strconv"
	"strings"

	"github.com/knavdeep152002/stream-sense/internal/models"
)

func (r *RoomHandler) createRoomInDB(room *models.Room) (*models.Room, error) {
	err := r.DB.Create(room).Error
	if err != nil {
		return nil, err
	}
	err = r.DB.Model(room).Association("Users").Append(room.AdminID)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rh *RoomHandler) getOrCreateRoomClient(id, videoId string, adminId uint) (rch *RoomClientHandler, err error) {
	return rh.roomClientHandlerPool.GetOrCreate(id, func() (*RoomClientHandler, error) {
		return NewRuntimeClientHandler(videoId, adminId), nil
	})
}

func (r *RoomHandler) getRoomAndVideoID(roomLink string) (room models.Room, err error) {
	roomLink = strings.TrimPrefix(roomLink, "/")
	fullPath := strings.Split(roomLink, "@")
	if len(fullPath) != 2 {
		err = errors.New("Invalid room link")
		return
	}
	id, err := strconv.Atoi(fullPath[0])
	if err != nil {
		return
	}
	roomId := uint(id)
	videoId := fullPath[1]
	// check if room exists and is active
	err = r.DB.Where("id=? AND video_id=? AND active=?", roomId, videoId, true).First(&room).Error
	if err != nil {
		return
	}
	return
}

func (r *RoomHandler) checkUserIsAdmin(roomId uint, userId uint) (bool, error) {
	var room models.Room
	err := r.DB.Where("id=?", roomId).First(&room).Error
	if err != nil {
		return false, err
	}
	return room.AdminID == userId, nil
}

func (r *RoomHandler) addUserToRoom(roomId uint, userId uint) error {
	var room models.Room
	err := r.DB.Where("id=?", roomId).First(&room).Error
	if err != nil {
		return err
	}
	// check if user is already in the room
	present := false
	for _, user := range room.Users {
		if user.ID == userId {
			present = true
			break
		}
	}
	if !present {
		err = r.DB.Model(&room).Association("Users").Append(userId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RoomHandler) removeUserFromRoom(roomId uint, userId uint) {
	var room models.Room
	err := r.DB.Where("id=?", roomId).First(&room).Error
	if err != nil {
		return
	}
	// check if user is already in the room
	present := false
	for _, user := range room.Users {
		if user.ID == userId {
			present = true
			break
		}
	}
	if present {
		err = r.DB.Model(&room).Association("Users").Delete(userId)
		if err != nil {
			return
		}
	}
}
