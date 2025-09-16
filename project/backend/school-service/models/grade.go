package models

import (
	"time"

	"gorm.io/gorm"
)

type Grade struct {
	gorm.Model
	Value   int       `gorm:"not null"`
	Date    time.Time `gorm:"not null"`
	Student Student   `gorm:"foreignKey:StudentID"`
	Subject Subject   `gorm:"foreignKey:SubjectID"`
	Teacher Teacher   `gorm:"foreignKey:TeacherID"`
}
