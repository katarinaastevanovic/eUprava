package models

import "gorm.io/gorm"

type Subject struct {
	gorm.Model
	Name     string    `gorm:"size:100;not null"`
	Teachers []Teacher `gorm:"many2many:subject_teachers"`
	Students []Student `gorm:"many2many:subject_students"`
}
