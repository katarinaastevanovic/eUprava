package services

import (
	"medical-service/database"
	"medical-service/models"
)

func CreateMedicalCertificate(cert *models.MedicalCertificate) error {
	var req models.Request
	if err := database.DB.First(&req, cert.RequestId).Error; err != nil {
		return err
	}

	var record models.MedicalRecord
	if err := database.DB.First(&record, req.MedicalRecordId).Error; err != nil {
		return err
	}

	cert.PatientId = record.PatientId

	return database.DB.Create(cert).Error
}
