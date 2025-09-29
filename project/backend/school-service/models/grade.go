package models

import (
	"time"

	"gorm.io/gorm"
)

type Grade struct {
	gorm.Model
	Value int       `gorm:"not null"`
	Date  time.Time `gorm:"not null"`

	StudentID uint    `gorm:"not null"`
	Student   Student `gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	SubjectID uint    `gorm:"not null"`
	Subject   Subject `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	TeacherID uint    `gorm:"not null"`
	Teacher   Teacher `gorm:"foreignKey:TeacherID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
