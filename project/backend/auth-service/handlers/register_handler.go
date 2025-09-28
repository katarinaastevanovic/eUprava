package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"auth-service/services"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req services.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := services.RegisterUser(req)
	if err != nil {
		switch err {
		case services.ErrUsernameExists:
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		case services.ErrInvalidUsername:
			http.Error(w, "Username must be at least 4 characters long", http.StatusBadRequest)
			return
		case services.ErrUMCNExists:
			http.Error(w, "UMCN already exists", http.StatusConflict)
			return
		case services.ErrEmailExists:
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		case services.ErrInvalidUMCN:
			http.Error(w, "UMCN must be exactly 13 characters", http.StatusBadRequest)
			return
		case services.ErrInvalidPassword:
			http.Error(w, "Password must be at least 8 characters long, contain one uppercase, one lowercase, and one special character", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
	}

	medicalReq := map[string]interface{}{
		"userId": user.ID,
	}

	body, _ := json.Marshal(medicalReq)

	resp, err := http.Post("http://medical-service:8082/medical-records", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusCreated {
		http.Error(w, "User registered but failed to create medical record", http.StatusAccepted)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully and medical record created"))
}

func CheckUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	available, err := services.CheckUsernameAvailable(username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !available {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
