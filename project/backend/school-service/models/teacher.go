package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	UserID      uint      `gorm:"not null"`
	Subjects    []Subject `gorm:"many2many:subject_teachers;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Title       string    `gorm:"size:100;not null"`
	HeadTeacher bool      `gorm:"default:false"`
}
