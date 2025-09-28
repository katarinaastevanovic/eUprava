package models

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	UserId       uint                 `gorm:"not null"`
	Certificates []MedicalCertificate `gorm:"foreignKey:PatientId"`
	DoctorID     *uint                `gorm:"column:doctor_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
