package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	docs "github.com/knavdeep152002/stream-sense/docs"
	"github.com/knavdeep152002/stream-sense/internal/constants"
	streamsense "github.com/knavdeep152002/stream-sense/internal/streamsense"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Length", "Content-Type", "Cookie"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
	}))
	pathPrefix := constants.PATH_PREFIX
	docs.SwaggerInfo.BasePath = pathPrefix + "/api/v1"

	ss := streamsense.NewStreamSense()
	ss.RegisterGroup(r)

	r.GET(pathPrefix+"/api/v1/doc/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(":8000")
}
