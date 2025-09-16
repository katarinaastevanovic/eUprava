package models

import "gorm.io/gorm"

type MedicalRecord struct {
	gorm.Model
	PatientID       uint   `gorm:"not null"`
	Allergies       string `gorm:"type:text"`
	ChronicDiseases string `gorm:"type:text"`
	LastUpdate      string
	Examinations    []Examination `gorm:"foreignKey:MedicalRecordID"`
	Requests        []Request     `gorm:"foreignKey:MedicalRecordID"`
}
