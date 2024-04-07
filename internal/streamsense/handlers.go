package streamsense

import (
	"github.com/gin-gonic/gin"
	"github.com/knavdeep152002/stream-sense/internal/constants"
	"github.com/knavdeep152002/stream-sense/internal/fs"
	"github.com/knavdeep152002/stream-sense/internal/redis"
)

type StreamSense struct{}

func (s *StreamSense) RegisterGroup(r *gin.Engine) {
	pathPrefix := constants.PATH_PREFIX

	ssGroup := r.Group(pathPrefix + "/api/v1")
	ssGroup.POST("/upload", fs.UploadChunk)
	ssGroup.POST("/complete", fs.CompleteUpload)

}

func NewStreamSense() *StreamSense {
	go redis.Observe()
	return &StreamSense{}
}
