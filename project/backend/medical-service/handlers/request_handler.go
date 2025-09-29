package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"medical-service/models"
	"medical-service/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func getUserIdFromJWT(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("no Authorization header")
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return 0, errors.New("invalid Authorization header")
	}

	token := parts[1]
	payloadPart := strings.Split(token, ".")
	if len(payloadPart) < 2 {
		return 0, errors.New("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(payloadPart[1])
	if err != nil {
		return 0, err
	}

	var decoded struct {
		Sub uint `json:"sub"`
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return 0, err
	}

	return decoded.Sub, nil
}

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	record, err := services.GetMedicalRecordByUserId(userId)
	if err != nil {
		http.Error(w, "Medical record not found", http.StatusBadRequest)
		return
	}

	var req struct {
		DoctorId uint                     `json:"doctorId"`
		Type     models.TypeOfExamination `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newReq := models.Request{
		MedicalRecordId: record.ID,
		DoctorId:        req.DoctorId,
		Type:            req.Type,
		Status:          models.REQUESTED,
	}

	if err := services.CreateRequest(&newReq); err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newReq)
}

func GetRequestsByPatient(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requests, err := services.GetRequestsByPatientUser(userId)
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
