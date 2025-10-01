package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
