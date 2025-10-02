package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"school-service/services"
	"sync"
	"time"
)

var userRequests = make(map[uint][]time.Time)
var mu sync.Mutex

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

	userId, err := services.GetUserIdFromJWT(jwtToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mu.Lock()
	now := time.Now()
	windowStart := now.Add(-time.Minute)
	timestamps := userRequests[userId]

	var updated []time.Time
	for _, t := range timestamps {
		if t.After(windowStart) {
			updated = append(updated, t)
		}
	}

	if len(updated) >= 5 {
		mu.Unlock()
		http.Error(w, "Too many requests, try again later", http.StatusTooManyRequests)
		return
	}

	updated = append(updated, now)
	userRequests[userId] = updated
	mu.Unlock()

	err = services.CreateRequest(req, jwtToken)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	log.Printf("[CreateRequest] Request successfully created for user %d", userId)
	w.WriteHeader(http.StatusCreated)
}
