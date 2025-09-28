package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	SubjectID uint   `gorm:"not null"`
	Title     string `gorm:"size:100;not null"`
}
