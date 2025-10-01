package services

import (
	"encoding/json"
	"fmt"
	"medical-service/database"
	"medical-service/models"
	"net/http"
	"strings"
)

func CreateRequest(req *models.Request) error {
	req.Status = models.REQUESTED
	if err := database.DB.Create(req).Error; err != nil {
		return err
	}

	var record models.MedicalRecord
	if err := database.DB.First(&record, req.MedicalRecordId).Error; err != nil {
		fmt.Println("Failed to get medical record:", err)
	}

	var patient models.Patient
	if err := database.DB.First(&patient, record.PatientId).Error; err != nil {
		fmt.Println("Failed to get patient:", err)
	}

	studentName := getStudentNameFromAuth(patient.UserId)

	message := fmt.Sprintf("You have a new examination request from %s. Type: %s.", studentName, req.Type)

	if err := CreateNotification(req.DoctorId, message); err != nil {
		fmt.Println("Failed to create notification:", err)
	}

	return nil
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

func GetRequestsByDoctorWithStudentPaginated(doctorId uint, page, pageSize int, status models.TypeOfRequest, search string) ([]RequestWithStudent, int, error) {
	var requests []models.Request
	query := database.DB.Where("doctor_id = ?", doctorId)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	var filtered []RequestWithStudent
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

		if search != "" && !strings.Contains(strings.ToLower(studentName), strings.ToLower(search)) {
			continue
		}

		filtered = append(filtered, RequestWithStudent{
			ID:              req.ID,
			MedicalRecordId: req.MedicalRecordId,
			DoctorId:        req.DoctorId,
			Type:            req.Type,
			Status:          req.Status,
			StudentName:     studentName,
		})
	}

	total := len(filtered)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paginated := filtered[start:end]

	totalPages := (total + pageSize - 1) / pageSize
	return paginated, totalPages, nil
}

func GetRequestsByDoctorWithStudentPaginatedCustomFilters(
	doctorId uint,
	page, pageSize int,
	status string,
	search string,
	reqType string,
	sortPending bool, // flag za globalno sortiranje
) ([]RequestWithStudent, int, error) {

	var requests []models.Request
	query := database.DB.Where("doctor_id = ?", doctorId)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if reqType != "" {
		query = query.Where("type = ?", reqType)
	}

	if err := query.Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	var filtered []RequestWithStudent
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

		if search != "" && !strings.Contains(strings.ToLower(studentName), strings.ToLower(search)) {
			continue
		}

		filtered = append(filtered, RequestWithStudent{
			ID:              req.ID,
			MedicalRecordId: req.MedicalRecordId,
			DoctorId:        req.DoctorId,
			Type:            req.Type,
			Status:          req.Status,
			StudentName:     studentName,
		})
	}

	if sortPending {
		filtered = SortRequestsPendingFirst(filtered)
	}

	total := len(filtered)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paginated := filtered[start:end]

	totalPages := (total + pageSize - 1) / pageSize
	return paginated, totalPages, nil
}

func SortRequestsPendingFirst(requests []RequestWithStudent) []RequestWithStudent {
	sorted := make([]RequestWithStudent, len(requests))
	copy(sorted, requests)

	pending := []RequestWithStudent{}
	others := []RequestWithStudent{}

	for _, r := range sorted {
		if r.Status == "REQUESTED" {
			pending = append(pending, r)
		} else {
			others = append(others, r)
		}
	}

	return append(pending, others...)
}
