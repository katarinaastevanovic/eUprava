package services

import (
	"encoding/json"
	"fmt"
	"medical-service/database"
	"medical-service/models"
	"net/http"
	"time"
)

type AuthUser struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"lastName"`
	UMCN      string `json:"umcn"`
	BirthDate string `json:"birthDate"`
	Gender    string `json:"gender"`
}

func CreateMedicalRecord(patientID uint) (*models.MedicalRecord, error) {
	record := models.MedicalRecord{
		PatientId:       patientID,
		Allergies:       "",
		ChronicDiseases: "",
		LastUpdate:      time.Now().Format("2006-01-02"),
	}

	if err := database.DB.Create(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func GetMedicalRecordByUserID(userID uint) (*models.MedicalRecord, error) {
	var record models.MedicalRecord
	if err := database.DB.Where("patient_id = ?", userID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func UpdateMedicalRecord(userID uint, allergies, chronicDiseases string) (*models.MedicalRecord, error) {
	var record models.MedicalRecord
	if err := database.DB.Where("patient_id = ?", userID).First(&record).Error; err != nil {
		return nil, err
	}

	record.Allergies = allergies
	record.ChronicDiseases = chronicDiseases
	record.LastUpdate = time.Now().Format("2006-01-02")

	if err := database.DB.Save(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func GetPatientFromAuth(userID uint) (*AuthUser, error) {
	url := fmt.Sprintf("http://authservice:8081/patients/%d", userID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	var user AuthUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &user, nil
}

func GetPatientByUserID(userID uint) (*models.Patient, error) {
	var patient models.Patient
	if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func GetMedicalRecordByPatientID(patientID uint) (*models.MedicalRecord, error) {
	var record models.MedicalRecord
	if err := database.DB.Where("patient_id = ?", patientID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func GetMedicalRecordIdByRequest(requestId uint) (uint, error) {
	var request models.Request
	err := database.DB.First(&request, requestId).Error
	if err != nil {
		return 0, err
	}
	return request.MedicalRecordId, nil
}
