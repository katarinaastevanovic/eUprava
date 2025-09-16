package models

import (
	"time"

	"gorm.io/gorm"
)

type Absence struct {
	gorm.Model
	Type    AbsenceType `gorm:"type:varchar(20);not null"`
	Date    time.Time   `gorm:"not null"`
	Student Student     `gorm:"foreignKey:StudentID"`
	//StudentID   uint
	Subject Subject `gorm:"foreignKey:SubjectID"`
}
