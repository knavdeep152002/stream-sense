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
	"gorm.io/gorm"
)

type StreamSense struct {
	authModule *auth.Auth
	authCheck  *middlewares.AuthCheck
	DB         *gorm.DB
}

func (s *StreamSense) RegisterGroup(r *gin.Engine) {
	pathPrefix := constants.PATH_PREFIX

	ssGroup := r.Group(pathPrefix + "/api/v1")
	ssGroup.POST("/upload", fs.UploadChunk)
	ssGroup.POST("/complete", fs.CompleteUpload)
	ssGroup.POST("/qa/:videoId", openai.VideoIntelligence)
	ssGroup.POST("/auth/register", s.authModule.CreateUser)
	ssGroup.POST("/auth/login", s.authModule.Login)
}

func NewStreamSense() *StreamSense {
	go redis.Observe()
	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	return &StreamSense{
		authModule: auth.CreateAuth(db),
		DB:         db,
		authCheck:  middlewares.CreateAuthCheck(db),
	}
}
