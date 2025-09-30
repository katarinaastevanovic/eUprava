package handlers

import (
	"encoding/json"
	"log"
	"medical-service/models"
	"medical-service/services"
	"net/http"
)

func CreatePatientHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserId   uint   `json:"userId"`
		Name     string `json:"name"`
		LastName string `json:"lastName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	patient := models.Patient{
		UserId: req.UserId,
	}

	createdPatient, err := services.CreatePatient(patient)
	if err != nil {
		log.Printf("Failed to create patient: %v", err)
		http.Error(w, "Failed to create patient", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPatient)
}
