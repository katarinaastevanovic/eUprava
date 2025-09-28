package services

import (
	"medical-service/database"
	"medical-service/models"
)

func CreatePatient(patient models.Patient) (models.Patient, error) {
	if err := database.DB.Create(&patient).Error; err != nil {
		return models.Patient{}, err
	}
	return patient, nil
}
