package services

import (
	"encoding/json"
	"fmt"
	"medical-service/database"
	"medical-service/models"
	"net/http"
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

type RequestWithStudent struct {
	ID              uint                     `json:"id"`
	MedicalRecordId uint                     `json:"medicalRecordId"`
	DoctorId        uint                     `json:"doctorId"`
	Type            models.TypeOfExamination `json:"type"`
	Status          models.TypeOfRequest     `json:"status"`
	StudentName     string                   `json:"studentName"`
}

func GetRequestsByDoctorWithStudent(doctorId uint) ([]RequestWithStudent, error) {
	var requests []models.Request
	if err := database.DB.Where("doctor_id = ?", doctorId).Find(&requests).Error; err != nil {
		return nil, err
	}

	var results []RequestWithStudent
	for _, req := range requests {
		var record models.MedicalRecord
		if err := database.DB.First(&record, req.MedicalRecordId).Error; err != nil {
			continue
		}

		var patient models.Patient
		if err := database.DB.First(&patient, record.PatientId).Error; err != nil {
			continue
		}

		studentName := getStudentNameFromAuth(patient.UserId)

		results = append(results, RequestWithStudent{
			ID:              req.ID,
			MedicalRecordId: req.MedicalRecordId,
			DoctorId:        req.DoctorId,
			Type:            req.Type,
			Status:          req.Status,
			StudentName:     studentName,
		})
	}

	return results, nil
}

func getStudentNameFromAuth(userId uint) string {
	url := fmt.Sprintf("http://auth-service:8081/users/%d", userId)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("ID: %d", userId)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("ID: %d", userId)
	}

	var user struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		LastName string `json:"lastName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return fmt.Sprintf("ID: %d", userId)
	}

	return fmt.Sprintf("%s %s", user.Name, user.LastName)
}

func GetApprovedRequestsByDoctorWithStudent(doctorId uint) ([]RequestWithStudent, error) {
	var requests []models.Request
	if err := database.DB.Where("doctor_id = ? AND status = ?", doctorId, models.APPROVED).Find(&requests).Error; err != nil {
		return nil, err
	}

	var results []RequestWithStudent
	for _, req := range requests {
		var record models.MedicalRecord
		if err := database.DB.First(&record, req.MedicalRecordId).Error; err != nil {
			continue
		}

		var patient models.Patient
		if err := database.DB.First(&patient, record.PatientId).Error; err != nil {
			continue
		}

		studentName := getStudentNameFromAuth(patient.UserId)

		results = append(results, RequestWithStudent{
			ID:              req.ID,
			MedicalRecordId: req.MedicalRecordId,
			DoctorId:        req.DoctorId,
			Type:            req.Type,
			Status:          req.Status,
			StudentName:     studentName,
		})
	}

	return results, nil
}

func GetRequestById(requestId uint) (*models.Request, error) {
	var req models.Request
	if err := database.DB.First(&req, requestId).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func GetRequestsByDoctorWithStudentPaginated(doctorId uint, page, pageSize int, status models.TypeOfRequest) ([]RequestWithStudent, int, error) {
	var total int64
	query := database.DB.Model(&models.Request{}).Where("doctor_id = ?", doctorId)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var requests []models.Request
	query = database.DB.Where("doctor_id = ?", doctorId)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Limit(pageSize).Offset(offset).Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	var results []RequestWithStudent
	for _, req := range requests {
		var record models.MedicalRecord
		if err := database.DB.First(&record, req.MedicalRecordId).Error; err != nil {
			fmt.Printf("MedicalRecord not found for Request ID %d, MedicalRecordId %d\n", req.ID, req.MedicalRecordId)
			continue
		}

		var patient models.Patient
		if err := database.DB.First(&patient, record.PatientId).Error; err != nil {
			fmt.Printf("Patient not found for MedicalRecord ID %d, PatientId %d\n", record.ID, record.PatientId)
			continue
		}

		fmt.Printf("Found Patient: ID=%d, UserId=%d\n", patient.ID, patient.UserId)

		studentName := getStudentNameFromAuth(patient.UserId)
		fmt.Printf("Student name resolved: %s\n", studentName)

		results = append(results, RequestWithStudent{
			ID:              req.ID,
			MedicalRecordId: req.MedicalRecordId,
			DoctorId:        req.DoctorId,
			Type:            req.Type,
			Status:          req.Status,
			StudentName:     studentName,
		})
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return results, totalPages, nil
}
