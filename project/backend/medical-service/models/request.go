package models

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	MedicalRecordId        uint              `gorm:"not null" json:"medicalRecordId"`
	DoctorId               uint              `gorm:"not null" json:"doctorId"`
	Type                   TypeOfExamination `gorm:"type:varchar(20);not null" json:"type"`
	Status                 TypeOfRequest     `gorm:"type:varchar(20);not null" json:"status"`
	NeedMedicalCertificate *bool             `gorm:"" json:"needMedicalCertificate"`
	Examinations           []Examination     `gorm:"foreignKey:RequestId" json:"examinations"`
}
