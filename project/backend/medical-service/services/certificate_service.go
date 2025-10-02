package services

import (
	"errors"
	"medical-service/models"

	"gorm.io/gorm"
)

type CertificateService struct {
	db *gorm.DB
}

func NewCertificateService(db *gorm.DB) *CertificateService {
	return &CertificateService{db: db}
}

// Proverava da li user ima neki medical certificate
func (s *CertificateService) HasCertificate(userId uint) (bool, error) {
	var patient models.Patient
	err := s.db.Preload("Certificates").Where("user_id = ?", userId).First(&patient).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // nema pacijenta → nema ni sertifikat
	}
	if err != nil {
		return false, err
	}

	// ako lista sertifikata nije prazna → ima sertifikat
	return len(patient.Certificates) > 0, nil
}
