package models

import "gorm.io/gorm"

type MedicalCertificate struct {
	gorm.Model
	RequestId uint              `json:"requestId" gorm:"not null"`
	Request   Request           `gorm:"foreignKey:RequestId"`
	PatientId uint              `json:"patientId" gorm:"not null"`
	DoctorId  uint              `json:"doctorId" gorm:"not null"`
	Date      string            `json:"date" gorm:"not null"`
	Type      TypeOfCertificate `json:"type" gorm:"type:varchar(20);not null"`
	Note      string            `json:"note" gorm:"type:text"`
}
