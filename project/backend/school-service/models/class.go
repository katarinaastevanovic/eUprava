package models

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Title         string    `gorm:"size:100;not null"`
	Year          int       `gorm:"not null"`
	Subjects      []Subject `gorm:"many2many:class_subjects;"`
	HeadTeacherID uint      `gorm:"not null"`
	HeadTeacher   Teacher   `gorm:"foreignKey:HeadTeacherID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
