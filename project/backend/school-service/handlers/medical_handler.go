package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"school-service/services"
)

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req services.RequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	jwtToken := r.Header.Get("Authorization")
	if jwtToken == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if len(jwtToken) > 7 && jwtToken[:7] == "Bearer " {
		jwtToken = jwtToken[7:]
	}

	err := services.CreateRequest(req, jwtToken)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	log.Printf("[CreateRequest] Request successfully created")
	w.WriteHeader(http.StatusCreated)
}
