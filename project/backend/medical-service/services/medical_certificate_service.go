package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"medical-service/database"
	"medical-service/models"
	"net/http"
)

func CreateMedicalCertificate(cert *models.MedicalCertificate) error {
	var req models.Request
	if err := database.DB.First(&req, cert.RequestId).Error; err != nil {
		return err
	}

	var record models.MedicalRecord
	if err := database.DB.First(&record, req.MedicalRecordId).Error; err != nil {
		return err
	}

	cert.PatientId = record.PatientId

	if err := database.DB.Create(cert).Error; err != nil {
		return err
	}

	var patient models.Patient
	if err := database.DB.First(&patient, cert.PatientId).Error; err != nil {
		return err
	}

	message := fmt.Sprintf("New Medical Certificate created of type %s dated %s", cert.Type, cert.Date)
	if err := NotifyStudent(patient.UserId, message); err != nil {
		log.Printf("[NotifyStudent] Failed to notify student %d: %v", patient.UserId, err)
	} else {
		log.Printf("[NotifyStudent] Notification sent to student %d", patient.UserId)
	}

	return nil
}

func NotifyStudent(userId uint, message string) error {
	notif := map[string]interface{}{
		"studentId": userId,
		"message":   message,
	}

	body, _ := json.Marshal(notif)
	resp, err := http.Post("http://school-service:8083/notifications", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("school-service returned status %d", resp.StatusCode)
	}

	return nil
}
