package models

import "time"

type Notification struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	UserId    uint      `gorm:"not null" json:"userId"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Read      bool      `gorm:"default:false" json:"read"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
