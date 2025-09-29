package services

import (
	"medical-service/database"
	"medical-service/models"
)

func CreateRequest(req *models.Request) error {
	req.Status = models.REQUESTED
	return database.DB.Create(req).Error
}

func GetRequestsByPatientUser(userId uint) ([]models.Request, error) {
	record, err := GetMedicalRecordByUserId(userId)
	if err != nil {
		return nil, err
	}

	var requests []models.Request
	err = database.DB.Where("medical_record_id = ?", record.ID).Find(&requests).Error
	return requests, err
}

func GetMedicalRecordByUserId(userId uint) (*models.MedicalRecord, error) {
	var patient models.Patient
	if err := database.DB.Where("user_id = ?", userId).First(&patient).Error; err != nil {
		return nil, err
	}

	var record models.MedicalRecord
	if err := database.DB.Where("patient_id = ?", patient.ID).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func GetRequestsByDoctor(doctorId uint) ([]models.Request, error) {
	var requests []models.Request
	err := database.DB.Where("doctor_id = ?", doctorId).Find(&requests).Error
	return requests, err
}

func UpdateRequestStatus(id uint, status models.TypeOfRequest) (*models.Request, error) {
	var req models.Request
	if err := database.DB.First(&req, id).Error; err != nil {
		return nil, err
	}
	req.Status = status
	if err := database.DB.Save(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}
