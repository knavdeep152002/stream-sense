package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Active  bool   `json:"active"`
	VideoID string `json:"video_id"`
	Users   []User `gorm:"many2many:room_users;"`
	AdminID uint   `json:"admin_id"`
	Admin   User   `gorm:"foreignKey:AdminID;references:ID"`
}