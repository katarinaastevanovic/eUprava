package handlers

import (
	"encoding/json"
	"medical-service/models"
	"medical-service/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateExamination(w http.ResponseWriter, r *http.Request) {
	var exam models.Examination
	if err := json.NewDecoder(r.Body).Decode(&exam); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := services.CreateExamination(&exam); err != nil {
		http.Error(w, "Failed to create examination", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exam)
}

func GetExaminationByRequest(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["requestId"]
	requestId, _ := strconv.Atoi(idStr)

	exam, err := services.GetExaminationByRequestId(uint(requestId))
	if err != nil {
		http.Error(w, "Examination not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(exam)
}
