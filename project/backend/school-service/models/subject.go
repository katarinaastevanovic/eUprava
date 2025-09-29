package models

import "gorm.io/gorm"

type Subject struct {
	gorm.Model
	Name     string    `gorm:"size:100;not null"`
	Grades   []Grade   `gorm:"foreignKey:SubjectID"`
	Absences []Absence `gorm:"foreignKey:SubjectID"`
}
