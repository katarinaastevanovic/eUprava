package models

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	UMCN     string `gorm:"unique;not null"`
	Name     string `gorm:"size:100;not null"`
	LastName string `gorm:"size:100;not null"`
	Email    string `gorm:"size:100;unique;not null"`
	Username string `gorm:"size:50;unique;not null"`
	Password string `gorm:"size:255;not null"`
	Role     string `gorm:"size:50;not null"`
}
