package handlers

import (
	"encoding/json"
	"net/http"

	"auth-service/middleware"
	"auth-service/services"
)

type CurrentUserResponse struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	user, err := services.FindUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := CurrentUserResponse{
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
