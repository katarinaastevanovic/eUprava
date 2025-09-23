package services

import (
	"auth-service/database"
	"auth-service/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

func CreateUserFromFirebase(uid, email, username, name, lastName, umcn string) (*models.Member, error) {
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
