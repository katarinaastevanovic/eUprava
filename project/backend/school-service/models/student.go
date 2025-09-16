package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	UserID           uint `gorm:"not null;index"`
	NumberOfAbsences int  `gorm:"not null"`
	ClassID          uint `gorm:"not null"`
}
