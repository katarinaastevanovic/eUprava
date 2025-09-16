package models

import "gorm.io/gorm"

type Examination struct {
	gorm.Model
	RequestID       uint   `gorm:"not null"`
	MedicalRecordID uint   `gorm:"not null"`
	Diagnosis       string `gorm:"type:text"`
	Therapy         string `gorm:"type:text"`
	Note            string `gorm:"type:text"`
}
