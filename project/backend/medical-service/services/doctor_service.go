package services

import (
	"medical-service/database"
	"medical-service/models"
)

func CreateDoctor(doctor models.Doctor) (models.Doctor, error) {
	if err := database.DB.Create(&doctor).Error; err != nil {
		return models.Doctor{}, err
	}
	return doctor, nil
}
