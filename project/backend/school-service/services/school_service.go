package services

import (
	"errors"
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

func (s *SchoolService) CreateAbsence(absence *models.Absence) error {
	if absence.StudentID == 0 {
		return errors.New("studentId is required")
	}
	if absence.SubjectID == 0 {
		return errors.New("subjectId is required")
	}

	if absence.Type == "" {
		absence.Type = models.Pending
	}

	if absence.Date.IsZero() {
		absence.Date = time.Now()
	}

	return s.DB.Create(absence).Error
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

func (s *SchoolService) GetStudentsByClassID(classID uint) ([]models.Student, error) {
	var students []models.Student

	err := s.DB.Preload("Absences").Preload("Grades").Preload("Class").
		Where("class_id = ?", classID).Find(&students).Error
	if err != nil {
		return nil, err
	}

	return students, nil
}
