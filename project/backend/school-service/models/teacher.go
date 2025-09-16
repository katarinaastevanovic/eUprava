package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Subjects    []Subject `gorm:"many2many:teacher_subjects;"`
	Title       string    `gorm:"size:100;not null"`
	HeadTeacher bool      `gorm:"default:false"`
}
