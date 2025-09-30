package services

import (
	"auth-service/database"
	"auth-service/models"
	"errors"
	"fmt"
	"strconv"

	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
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

func RegisterUser(req RegisterRequest) (*models.Member, error) {
	if len(req.UMCN) != 13 {
		return nil, ErrInvalidUMCN
	}

	if len(req.Username) < 4 {
		return nil, ErrInvalidUsername
	}

	if !validatePassword(req.Password) {
		return nil, ErrInvalidPassword
	}

	var count int64
	database.DB.Model(&models.Member{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, ErrUsernameExists
	}

	database.DB.Model(&models.Member{}).Where("umcn = ?", req.UMCN).Count(&count)
	if count > 0 {
		return nil, ErrUMCNExists
	}

	database.DB.Model(&models.Member{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	birthDate, gender, err := ParseUMCN(req.UMCN)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &member, nil
}

func CheckUsernameAvailable(username string) (bool, error) {
	var count int64
	if err := database.DB.Model(&models.Member{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func ParseUMCN(umcn string) (time.Time, string, error) {
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

var JwtKey = []byte("super-secret-key")

func GenerateJWT(userID uint, username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"role":     role,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func FindUserByEmail(email string) (*models.Member, error) {
	var user models.Member
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByID(id uint) (*models.Member, error) {
	var user models.Member
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUserFromFirebase(uid *string, email, username, name, lastName, umcn string) (*models.Member, error) {
	user := models.Member{
		FirebaseUID: uid,
		Email:       email,
		Username:    username,
		Name:        name,
		LastName:    lastName,
		UMCN:        umcn,
		Role:        "user",
		Password:    "",
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByFirebaseUID(uid string) (*models.Member, error) {
	var user models.Member
	if err := database.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func GetAllDoctors() ([]models.Member, error) {
	var doctors []models.Member
	err := database.DB.Where("role = ?", "DOCTOR").Find(&doctors).Error
	return doctors, err
}

func GetAllStudents() ([]models.Member, error) {
	var students []models.Member
	if err := database.DB.Where("role = ?", "STUDENT").Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

type MembersBatchRequest struct {
	IDs []uint `json:"ids"`
}

type MemberResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

func GetMembersByIDs(ids []uint) ([]MemberResponse, error) {
	var members []models.Member
	if err := database.DB.Select("id, name, last_name, email").
		Where("id IN ?", ids).
		Find(&members).Error; err != nil {
		return nil, err
	}

	result := make([]MemberResponse, len(members))
	for i, m := range members {
		result[i] = MemberResponse{
			ID:       m.ID,
			Name:     m.Name,
			LastName: m.LastName,
			Email:    m.Email,
		}
	}

	return result, nil

}
