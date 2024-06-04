package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string    `json:"username" gorm:"unique"`
	Password string    `json:"password"`
	Uploads  []Uploads `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID;references:ID"` // One to many relationship
}
