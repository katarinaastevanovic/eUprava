package handlers

import (
	"encoding/json"
	"net/http"

	"auth-service/database"
	"auth-service/models"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	UMCN     string `json:"umcn"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash lozinke
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	member := models.Member{
		UMCN:     req.UMCN,
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     models.Role(req.Role),
	}

	if err := database.DB.Create(&member).Error; err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
func CheckUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var count int64
	if err := database.DB.Model(&models.Member{}).Where("username = ?", username).Count(&count).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
