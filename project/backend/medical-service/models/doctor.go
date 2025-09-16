package models

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	UserID   uint      `gorm:"not null"`
	Type     string    `gorm:"size:50;not null"`
	Patients []Patient `gorm:"foreignKey:DoctorID"`
}
