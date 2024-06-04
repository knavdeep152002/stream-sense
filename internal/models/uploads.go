package models

import (
	"gorm.io/gorm"
)

type Uploads struct {
	gorm.Model
	FileName string `json:"filename"`
	VideoId  string `json:"video_id"`
	UserID   uint   `json:"user_id"`
	User     User   `gorm:"foreignKey:UserID;references:ID"`
}
