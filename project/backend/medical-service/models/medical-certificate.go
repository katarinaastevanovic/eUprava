package models

import "gorm.io/gorm"

type MedicalCertificate struct {
	gorm.Model
	PatientID uint              `gorm:"not null"`
	DoctorID  uint              `gorm:"not null"`
	Date      string            `gorm:"not null"`
	Type      TypeOfCertificate `gorm:"type:varchar(20);not null"`
	Note      string            `gorm:"type:text"`
}
