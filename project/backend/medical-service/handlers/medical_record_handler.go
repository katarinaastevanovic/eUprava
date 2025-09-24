package handlers

import (
	"encoding/json"
	"medical-service/database"
	"medical-service/models"
	"net/http"
	"time"
)

type CreateMedicalRecordRequest struct {
	UserID uint `json:"userId"`
}

func CreateMedicalRecord(w http.ResponseWriter, r *http.Request) {
	var req CreateMedicalRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	record := models.MedicalRecord{
		PatientID:       req.UserID,
		Allergies:       "",
		ChronicDiseases: "",
		LastUpdate:      time.Now().Format("2006-01-02"),
	}

	if err := database.DB.Create(&record).Error; err != nil {
		http.Error(w, "Failed to create record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(record)
}
