package handlers

import (
	"auth-service/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	doctors, err := services.GetAllDoctors()
	if err != nil {
		http.Error(w, "Failed to fetch doctors", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(doctors)
}

func GetAllStudentsHandler(w http.ResponseWriter, r *http.Request) {
	students, err := services.GetAllStudents()
	if err != nil {
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["userId"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := services.FindUserByID(uint(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
