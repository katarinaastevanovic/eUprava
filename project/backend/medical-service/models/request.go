package models

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	MedicalRecordId uint              `gorm:"not null"`
	DoctorId        uint              `gorm:"not null"`
	Type            TypeOfExamination `gorm:"type:varchar(20);not null"`
	Status          TypeOfRequest     `gorm:"type:varchar(20);not null"`
	Examinations    []Examination     `gorm:"foreignKey:RequestId"`
}
