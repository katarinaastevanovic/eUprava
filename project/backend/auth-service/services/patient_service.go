package services

import (
	"auth-service/database"
	"auth-service/models"
	"time"
)

type AuthUser struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastName"`
	UMCN      string    `json:"umcn"`
	BirthDate time.Time `json:"birthDate"`
	Gender    string    `json:"gender"`
}

func GetPatientData(userID uint) (*AuthUser, error) {
	var member models.Member
	if err := database.DB.First(&member, userID).Error; err != nil {
		return nil, err
	}

	user := &AuthUser{
		ID:        member.ID,
		Name:      member.Name,
		LastName:  member.LastName,
		UMCN:      member.UMCN,
		BirthDate: member.BirthDate,
		Gender:    member.Gender,
	}

	return user, nil
}
