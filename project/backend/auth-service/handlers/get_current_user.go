package handlers

import (
	"encoding/json"
	"net/http"

	"auth-service/middleware"
	"auth-service/services"
)

type CurrentUserResponse struct {
	UserID    uint   `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Role      string `json:"role"`
	UMCN      string `json:"umcn,omitempty"`
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := services.FindUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := CurrentUserResponse{
		UserID:    user.ID,
		Username:  user.Username,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		BirthDate: user.BirthDate.Format("2006-01-02"),
		Gender:    string(user.Gender),
		Role:      string(user.Role),
		UMCN:      user.UMCN,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
