package models

import "gorm.io/gorm"

type Examination struct {
	gorm.Model
	RequestId       uint   `gorm:"not null" json:"requestId"`
	MedicalRecordId uint   `gorm:"not null" json:"medicalRecordId"`
	Diagnosis       string `gorm:"type:text" json:"diagnosis"`
	Therapy         string `gorm:"type:text" json:"therapy"`
	Note            string `gorm:"type:text" json:"note"`
}
