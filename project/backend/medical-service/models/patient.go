package models

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	UserId         uint                 `gorm:"not null"`
	MedicalRecords []MedicalRecord      `gorm:"foreignKey:PatientId"`
	Certificates   []MedicalCertificate `gorm:"foreignKey:PatientId"`
	DoctorID       uint
}
