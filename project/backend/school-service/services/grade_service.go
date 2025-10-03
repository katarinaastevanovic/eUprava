// services/grade_service.go
package services

import (
	"time"

	"school-service/models"

	"gorm.io/gorm"
)

type GradeService struct {
	db *gorm.DB
}

func NewGradeService(db *gorm.DB) *GradeService {
	return &GradeService{db: db}
}

func (s *GradeService) CreateGrade(value int, studentID, subjectID, teacherID uint) (*models.Grade, error) {
	grade := models.Grade{
		Value:     value,
		Date:      time.Now(),
		StudentID: studentID,
		SubjectID: subjectID,
		TeacherID: teacherID,
	}

	if err := s.db.Create(&grade).Error; err != nil {
		return nil, err
	}

	return &grade, nil
}
