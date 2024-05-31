package main

import (
	"github.com/knavdeep152002/stream-sense/internal/db"
	"github.com/knavdeep152002/stream-sense/internal/models"
	"github.com/knavdeep152002/stream-sense/internal/utils"
)

func init() {
	utils.LoadEnvs()
}

func main() {
	dbConn, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	err = dbConn.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
}
