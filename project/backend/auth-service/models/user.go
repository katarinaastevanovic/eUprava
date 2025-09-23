package models

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	UMCN        string `gorm:"unique;"`
	Name        string `gorm:"size:100;"`
	LastName    string `gorm:"size:100;"`
	Email       string `gorm:"size:100;unique;"`
	Username    string `gorm:"size:50;unique;"`
	Password    string `gorm:"size:255"`
	Role        Role   `gorm:"type:varchar(20);"`
	FirebaseUID string `gorm:"size:128;uniqueIndex"`
}
