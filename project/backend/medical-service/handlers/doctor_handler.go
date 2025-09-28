package handlers

import (
	"encoding/json"
	"log"
	"medical-service/models"
	"medical-service/services"
	"net/http"
)

func CreateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserId   uint   `json:"userId"`
		Name     string `json:"name"`
		LastName string `json:"lastName"`
		Type     string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	doctor := models.Doctor{
		UserId: req.UserId,
		Type:   req.Type,
	}

	createdDoctor, err := services.CreateDoctor(doctor)
	if err != nil {
		log.Printf("Failed to create doctor: %v", err)
		http.Error(w, "Failed to create doctor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdDoctor)
}
