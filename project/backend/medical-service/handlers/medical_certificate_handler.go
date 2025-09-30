package handlers

import (
	"encoding/json"
	"medical-service/models"
	"medical-service/services"
	"net/http"
)

func CreateMedicalCertificateHandler(w http.ResponseWriter, r *http.Request) {
	var cert models.MedicalCertificate
	if err := json.NewDecoder(r.Body).Decode(&cert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := services.CreateMedicalCertificate(&cert); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cert)
}
