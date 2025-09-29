package handlers

import (
	"auth-service/services"
	"encoding/json"
	"net/http"
)

type MembersBatchRequest struct {
	IDs []uint `json:"ids"`
}

type MemberResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

func GetMembersBatchHandler(w http.ResponseWriter, r *http.Request) {
	var req MembersBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		http.Error(w, "No IDs provided", http.StatusBadRequest)
		return
	}

	members, err := services.GetMembersByIDs(req.IDs)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}
