package services

import (
	"medical-service/database"
	"medical-service/models"
)

func CreateExamination(exam *models.Examination) error {
	return database.DB.Create(exam).Error
}

func GetExaminationByRequestId(requestId uint) (*models.Examination, error) {
	var exam models.Examination
	err := database.DB.Where("request_id = ?", requestId).First(&exam).Error
	if err != nil {
		return nil, err
	}
	return &exam, nil
}

func GetExaminationsByDoctor(doctorId uint) ([]models.Examination, error) {
	var exams []models.Examination
	err := database.DB.
		Joins("JOIN requests ON requests.id = examinations.request_id").
		Where("requests.doctor_id = ?", doctorId).
		Find(&exams).Error
	return exams, err
}
