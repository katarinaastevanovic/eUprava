package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte(getJwtKey())

func getJwtKey() string {
	if key := os.Getenv("JWT_SECRET"); key != "" {
		return key
	}
	return "super-secret-key"
}

func GetUserIdFromJWT(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	idFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("userID not found in token")
	}

	return uint(idFloat), nil
}

type RequestDTO struct {
	MedicalRecordId        uint   `json:"MedicalRecordId"`
	DoctorId               uint   `json:"DoctorId"`
	Type                   string `json:"Type"`
	Status                 string `json:"Status"`
	NeedMedicalCertificate *bool  `json:"NeedMedicalCertificate"`
}

func CreateRequest(req RequestDTO, jwtToken string) error {
	body, _ := json.Marshal(req)

	request, err := http.NewRequest("POST", "http://medical-service:8082/requests", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[CreateRequest Service] Failed to create HTTP request: %v", err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("medical-service returned status %d", resp.StatusCode)
	}

	return nil
}
