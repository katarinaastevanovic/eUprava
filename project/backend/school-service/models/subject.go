package models

import "gorm.io/gorm"

type Subject struct {
	gorm.Model
	Name     string    `gorm:"size:100;not null"`
	Teachers []Teacher `gorm:"many2many:subject_teachers;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Students []Student `gorm:"many2many:subject_students;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Grades   []Grade   `gorm:"foreignKey:SubjectID"`
	Absences []Absence `gorm:"foreignKey:SubjectID"`
}
