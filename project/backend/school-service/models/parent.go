package models

import "gorm.io/gorm"

type Parent struct {
	gorm.Model
	Children []Student `gorm:"many2many:parent_children;"`
}
