package streamsense

import (
	"github.com/gin-gonic/gin"
	"github.com/knavdeep152002/stream-sense/internal/auth"
	"github.com/knavdeep152002/stream-sense/internal/constants"
	"github.com/knavdeep152002/stream-sense/internal/db"
	"github.com/knavdeep152002/stream-sense/internal/fs"
	"github.com/knavdeep152002/stream-sense/internal/middlewares"
	"github.com/knavdeep152002/stream-sense/internal/openai"
	"github.com/knavdeep152002/stream-sense/internal/redis"
	"github.com/knavdeep152002/stream-sense/internal/room"
	"gorm.io/gorm"
)

type StreamSense struct {
	*auth.Auth
	*middlewares.AuthCheck
	*fs.FSHandler
	*room.RoomHandler
	DB *gorm.DB
}

func (s *StreamSense) RegisterGroup(r *gin.Engine) {
	pathPrefix := constants.PATH_PREFIX

	ssGroup := r.Group(pathPrefix + "/api/v1")
	// auth routes
	ssGroup.POST("/auth/register", s.CreateUser)
	ssGroup.POST("/auth/login", s.Login)

	// fs routes
	ssGroup.POST("/upload", s.Check, s.UploadChunk)
	ssGroup.POST("/complete", s.Check, s.CompleteUpload)
	ssGroup.GET("/uploads", s.Check, s.GetUserUploads)

	// room routes
	ssGroup.POST("/room/:videoId", s.Check, s.CreateRoom)
	ssGroup.GET("/rooms", s.Check, s.GetActiveRooms)
	ssGroup.GET("/room/*roomLink", s.Check, s.JoinRoom)
	ssGroup.GET("/serve/*roomLink", s.Check, s.ServeVideo)

	// openai routes
	ssGroup.POST("/qa/:videoId", s.Check, openai.VideoIntelligence)

}

func NewStreamSense() *StreamSense {
	go redis.Observe()
	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	return &StreamSense{
		DB:          db,
		Auth:        auth.CreateAuth(db),
		AuthCheck:   middlewares.CreateAuthCheck(db),
		FSHandler:   fs.CreateFSHandler(db),
		RoomHandler: room.NewRoomHandler(db),
	}
}
