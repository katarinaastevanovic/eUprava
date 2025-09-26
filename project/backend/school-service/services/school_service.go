package services

import (
	"school-service/models"

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
	// Pronađi studenta po user_id
	if err := s.DB.Where("user_id = ?", userID).First(&student).Error; err != nil {
		return nil, 0, err
	}

	var absences []models.Absence
	var count int64

	// Brojanje i učitavanje izostanaka
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
