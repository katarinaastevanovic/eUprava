package models

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	UserID         uint                 `gorm:"not null"`
	MedicalRecords []MedicalRecord      `gorm:"foreignKey:PatientID"`
	Certificates   []MedicalCertificate `gorm:"foreignKey:PatientID"`
	DoctorID       uint
}
