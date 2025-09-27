package models

import "gorm.io/gorm"

type MedicalRecord struct {
	gorm.Model
	PatientId       uint   `gorm:"not null"`
	Allergies       string `gorm:"type:text"`
	ChronicDiseases string `gorm:"type:text"`
	LastUpdate      string
	Examinations    []Examination `gorm:"foreignKey:MedicalRecordId"`
	Requests        []Request     `gorm:"foreignKey:MedicalRecordId"`
}
