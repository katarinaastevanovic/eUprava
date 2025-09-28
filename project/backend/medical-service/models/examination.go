package models

import "gorm.io/gorm"

type Examination struct {
	gorm.Model
	RequestId       uint   `gorm:"not null"`
	MedicalRecordId uint   `gorm:"not null"`
	Diagnosis       string `gorm:"type:text"`
	Therapy         string `gorm:"type:text"`
	Note            string `gorm:"type:text"`
}
