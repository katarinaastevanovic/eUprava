package models

import (
	"time"

	"gorm.io/gorm"
)

type Absence struct {
	gorm.Model
	Type      AbsenceType `gorm:"type:varchar(20);not null"`
	Date      time.Time   `gorm:"not null"`
	StudentID uint        `gorm:"not null"`
	Student   Student     `gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SubjectID uint        `gorm:"not null"`
	Subject   Subject     `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
