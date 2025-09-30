package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	UserID           uint      `gorm:"not null;index"`
	NumberOfAbsences int       `gorm:"not null"`
	ClassID          uint      `gorm:"not null"`
	Class            Class     `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Absences         []Absence `gorm:"foreignKey:StudentID"`
	Grades           []Grade   `gorm:"foreignKey:StudentID"`
}
