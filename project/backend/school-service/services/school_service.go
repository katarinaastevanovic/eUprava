package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"school-service/models"
	"time"

	"gorm.io/gorm"
)

type SchoolService struct {
	DB *gorm.DB
}

func NewSchoolService(db *gorm.DB) *SchoolService {
	return &SchoolService{DB: db}
}

type AbsenceDTO struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
	Date string `json:"date"`
}

type StudentDTO struct {
	ID               uint         `json:"id"`
	UserID           uint         `json:"userId"`
	Name             string       `json:"name"`
	LastName         string       `json:"lastName"`
	NumberOfAbsences int          `json:"numberOfAbsences"`
	Absences         []AbsenceDTO `json:"absences"`
}

type MemberResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
}

func (s *SchoolService) GetAbsenceCountForSubject(studentID uint, subjectID uint) (int64, error) {
	var count int64
	err := s.DB.Model(&models.Absence{}).
		Where("student_id = ? AND subject_id = ?", studentID, subjectID).
		Count(&count).Error
	return count, err
}

func (s *SchoolService) GetAbsencesByUserID(userID uint) ([]models.Absence, int64, error) {
	var student models.Student
	if err := s.DB.Where("user_id = ?", userID).First(&student).Error; err != nil {
		return nil, 0, err
	}

	var absences []models.Absence
	var count int64

	if err := s.DB.Model(&models.Absence{}).
		Where("student_id = ?", student.ID).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := s.DB.Preload("Subject").
		Where("student_id = ?", student.ID).
		Find(&absences).Error; err != nil {
		return nil, 0, err
	}

	return absences, count, nil
}

func (s *SchoolService) UpdateAbsenceType(absenceID uint, newType models.AbsenceType) error {
	return s.DB.Model(&models.Absence{}).
		Where("id = ?", absenceID).
		Update("type", newType).Error
}

func (s *SchoolService) CreateAbsences(absences []models.Absence) error {
	for i := range absences {
		if absences[i].StudentID == 0 {
			return errors.New("studentId is required")
		}
		if absences[i].SubjectID == 0 {
			return errors.New("subjectId is required")
		}

		if absences[i].Type == "" {
			absences[i].Type = models.Pending
		}

		if absences[i].Date.IsZero() {
			absences[i].Date = time.Now()
		}
	}

	return s.DB.Create(&absences).Error
}

func (s *SchoolService) GetClassesByTeacherID(teacherID uint) ([]models.Class, error) {
	var teacher models.Teacher
	if err := s.DB.First(&teacher, teacherID).Error; err != nil {
		return nil, err
	}

	var classes []models.Class
	err := s.DB.
		Joins("JOIN class_subjects cs ON cs.class_id = classes.id").
		Where("cs.subject_id = ?", teacher.SubjectID).
		Find(&classes).Error

	return classes, err
}

func (s *SchoolService) GetStudentsByClassID(classID uint) ([]StudentDTO, error) {
	var students []models.Student

	err := s.DB.Preload("Absences").
		Where("class_id = ?", classID).Find(&students).Error
	if err != nil {
		return nil, err
	}

	// ako nema učenika u razredu -> vrati prazan slice
	if len(students) == 0 {
		return []StudentDTO{}, nil
	}

	var userIDs []uint
	for _, st := range students {
		userIDs = append(userIDs, st.UserID)
	}

	// pokuša da uzme imena iz auth servisa
	memberMap, err := fetchMembersByIDs(userIDs)
	if err != nil {
		// ako auth ne radi, samo loguj i nastavi sa praznom mapom
		fmt.Println("warning: failed to fetch members:", err)
		memberMap = make(map[uint]MemberResponse)
	}

	var result []StudentDTO
	for _, st := range students {
		studentDTO := StudentDTO{
			ID:               st.ID,
			UserID:           st.UserID,
			NumberOfAbsences: st.NumberOfAbsences,
		}

		if m, ok := memberMap[st.UserID]; ok {
			studentDTO.Name = m.Name
			studentDTO.LastName = m.LastName
		}

		for _, abs := range st.Absences {
			studentDTO.Absences = append(studentDTO.Absences, AbsenceDTO{
				ID:   abs.ID,
				Type: string(abs.Type),
				Date: abs.Date.Format("2006-01-02"),
			})
		}

		result = append(result, studentDTO)
	}

	return result, nil
}

func fetchMembersByIDs(ids []uint) (map[uint]MemberResponse, error) {
	reqBody := map[string][]uint{"ids": ids}
	body, _ := json.Marshal(reqBody)

	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		return nil, fmt.Errorf("AUTH_SERVICE_URL nije postavljen")
	}

	endpoint := fmt.Sprintf("%s/api/members/batch", authURL)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth-service returned status %d", resp.StatusCode)
	}

	var members []MemberResponse
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, err
	}

	result := make(map[uint]MemberResponse)
	for _, m := range members {
		result[m.ID] = m
	}
	return result, nil
}

type StudentByUserDTO struct {
	ID     uint `json:"id"`
	UserID uint `json:"userId"`
}

func (s *SchoolService) GetStudentByUserID(userID uint) (*models.Student, error) {
	var student models.Student
	if err := s.DB.Where("user_id = ?", userID).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (s *SchoolService) GetStudentDTOByUserID(userID uint) (*StudentByUserDTO, error) {
	student, err := s.GetStudentByUserID(userID)
	if err != nil {
		return nil, err
	}

	dto := &StudentByUserDTO{
		ID:     student.ID,
		UserID: student.UserID,
	}

	return dto, nil
}

type FullStudentProfileDTO struct {
	ID               uint   `json:"id"`
	UserID           uint   `json:"user_id"`
	Name             string `json:"name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	NumberOfAbsences int    `json:"number_of_absences"`
	ClassID          uint   `json:"class_id"`
}

func (s *SchoolService) GetFullStudentProfileByUserID(userID uint) (*FullStudentProfileDTO, error) {
	student, err := s.GetStudentByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 1. pripremi JSON telo sa ID-jem
	bodyData, err := json.Marshal(map[string][]uint{
		"ids": {userID},
	})
	if err != nil {
		return nil, err
	}

	// 2. napravi POST zahtev ka batch ruti
	authURL := "http://auth-service:8080/api/members/batch"
	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, fmt.Errorf("failed to contact auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	// 3. dekodiraj listu korisnika
	var members []struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, fmt.Errorf("failed to decode auth response: %w", err)
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("no member found for userID %d", userID)
	}

	member := members[0]

	// 4. sastavi DTO
	dto := &FullStudentProfileDTO{
		ID:               student.ID,
		UserID:           student.UserID,
		Name:             member.Name,
		LastName:         member.LastName,
		Email:            member.Email,
		NumberOfAbsences: student.NumberOfAbsences,
		ClassID:          student.ClassID,
	}
	return dto, nil
}
