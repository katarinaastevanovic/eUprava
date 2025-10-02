package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserId  uint   `json:"userId" gorm:"not null"`
	Message string `json:"message" gorm:"not null"`
	Read    bool   `json:"read" gorm:"default:false"`
}
