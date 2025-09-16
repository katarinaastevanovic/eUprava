package models

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	MedicalRecordID uint              `gorm:"not null"`
	DoctorID        uint              `gorm:"not null"`
	Type            TypeOfExamination `gorm:"type:varchar(20);not null"`
	Status          TypeOfRequest     `gorm:"type:varchar(20);not null"`
	Examinations    []Examination     `gorm:"foreignKey:RequestID"`
}

//
