package models

import (
	"time"

	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	UMCN        string    `gorm:"unique;not null"`
	Name        string    `gorm:"size:100;not null"`
	LastName    string    `gorm:"size:100;not null"`
	Email       string    `gorm:"size:100;unique;not null"`
	Username    string    `gorm:"size:50;unique;not null"`
	Password    string    `gorm:"size:255;not null"`
	Role        Role      `gorm:"type:varchar(20);not null"`
	BirthDate   time.Time `gorm:"not null"`
	Gender      string    `gorm:"type:char(1);not null"`
	FirebaseUID string    `gorm:"size:128;uniqueIndex"`
}
