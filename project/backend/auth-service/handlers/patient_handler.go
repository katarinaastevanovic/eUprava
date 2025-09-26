package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"auth-service/services"

	"github.com/gorilla/mux"
)

func GetPatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	user, err := services.GetPatientData(uint(userID))
	if err != nil {
		http.Error(w, "Failed to fetch patient data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
