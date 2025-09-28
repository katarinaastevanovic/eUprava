package handlers

import (
	"encoding/json"
	"medical-service/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CreateMedicalRecordRequest struct {
	PatientID uint `json:"patientId"`
}

func CreateMedicalRecord(w http.ResponseWriter, r *http.Request) {
	var req CreateMedicalRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	record, err := services.CreateMedicalRecord(req.PatientID)
	if err != nil {
		http.Error(w, "Failed to create record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(record)
}

func GetMedicalRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	record, err := services.GetMedicalRecordByUserID(uint(userID))
	if err != nil {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func UpdateMedicalRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	var updateData struct {
		Allergies       string `json:"allergies"`
		ChronicDiseases string `json:"chronicDiseases"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	record, err := services.UpdateMedicalRecord(uint(userID), updateData.Allergies, updateData.ChronicDiseases)
	if err != nil {
		http.Error(w, "Failed to update record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func GetFullMedicalRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	patient, err := services.GetPatientByUserID(uint(userID))
	if err != nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	record, err := services.GetMedicalRecordByPatientID(patient.ID)
	if err != nil {
		http.Error(w, "Medical record not found", http.StatusNotFound)
		return
	}

	authUser, err := services.GetPatientFromAuth(uint(userID))
	if err != nil {
		http.Error(w, "Failed to fetch patient data from auth service", http.StatusInternalServerError)
		return
	}

	fullRecord := struct {
		PatientID       uint        `json:"patientId"`
		Name            string      `json:"name"`
		LastName        string      `json:"lastName"`
		JMBG            string      `json:"jmbg"`
		BirthDate       string      `json:"birthDate"`
		Gender          string      `json:"gender"`
		Allergies       string      `json:"allergies"`
		ChronicDiseases string      `json:"chronicDiseases"`
		LastUpdate      string      `json:"lastUpdate"`
		Examinations    interface{} `json:"examinations"`
		Requests        interface{} `json:"requests"`
	}{
		PatientID:       record.PatientId,
		Name:            authUser.Name,
		LastName:        authUser.LastName,
		JMBG:            authUser.UMCN,
		BirthDate:       authUser.BirthDate,
		Gender:          authUser.Gender,
		Allergies:       record.Allergies,
		ChronicDiseases: record.ChronicDiseases,
		LastUpdate:      record.LastUpdate,
		Examinations:    record.Examinations,
		Requests:        record.Requests,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullRecord)
}
