package handlers

import (
	"encoding/json"
	"medical-service/models"
	"medical-service/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := services.CreateRequest(&req); err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func GetRequestsByPatient(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	patientId, _ := strconv.Atoi(idStr)

	requests, err := services.GetRequestsByPatient(uint(patientId))
	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

func GetRequestsByDoctor(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	doctorId, _ := strconv.Atoi(idStr)

	requests, err := services.GetRequestsByDoctor(uint(doctorId))
	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

func ApproveRequest(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	req, err := services.UpdateRequestStatus(uint(id), models.APPROVED)
	if err != nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(req)
}

func RejectRequest(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	req, err := services.UpdateRequestStatus(uint(id), models.REJECTED)
	if err != nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(req)
}
