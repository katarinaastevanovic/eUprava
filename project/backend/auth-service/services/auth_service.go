package services

import (
	"auth-service/database"
	"auth-service/models"
	"errors"
	"fmt"
	"strconv"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	UMCN     string `json:"umcn"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var (
	ErrUsernameExists  = errors.New("username already exists")
	ErrUMCNExists      = errors.New("UMCN already exists")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidUMCN     = errors.New("UMCN must be exactly 13 characters")
	ErrInvalidPassword = errors.New("password must be at least 8 characters long, contain one uppercase letter, one lowercase letter, and one special character")
	ErrInvalidUsername = errors.New("username must be at least 4 characters long")
)

func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char), unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasSpecial
}

func RegisterUser(req RegisterRequest) error {

	if len(req.UMCN) != 13 {
		return ErrInvalidUMCN
	}

	if len(req.Username) < 4 {
		return ErrInvalidUsername
	}

	if !validatePassword(req.Password) {
		return ErrInvalidPassword
	}

	var count int64
	database.DB.Model(&models.Member{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return ErrUsernameExists
	}

	database.DB.Model(&models.Member{}).Where("umcn = ?", req.UMCN).Count(&count)
	if count > 0 {
		return ErrUMCNExists
	}

	database.DB.Model(&models.Member{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	birthDate, gender, err := parseUMCN(req.UMCN)
	if err != nil {
		return err
	}

	member := models.Member{
		UMCN:      req.UMCN,
		Name:      req.Name,
		LastName:  req.LastName,
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		Role:      models.Role(req.Role),
		BirthDate: birthDate,
		Gender:    gender,
	}

	if err := database.DB.Create(&member).Error; err != nil {
		return err
	}

	return nil
}

func CheckUsernameAvailable(username string) (bool, error) {
	var count int64
	if err := database.DB.Model(&models.Member{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func parseUMCN(umcn string) (time.Time, string, error) {
	if len(umcn) != 13 {
		return time.Time{}, "", ErrInvalidUMCN
	}

	day, _ := strconv.Atoi(umcn[0:2])
	month, _ := strconv.Atoi(umcn[2:4])
	yearPart := umcn[4:7]

	y, _ := strconv.Atoi(yearPart)
	var year int
	if y <= 24 {
		year = 2000 + y
	} else {
		year = 1000 + y
	}

	birthDate, err := time.Parse("02-01-2006",
		fmt.Sprintf("%02d-%02d-%d", day, month, year))
	if err != nil {
		return time.Time{}, "", err
	}

	seq, _ := strconv.Atoi(umcn[7:10])
	gender := "M"
	if seq >= 500 {
		gender = "F"
	}

	return birthDate, gender, nil
}
