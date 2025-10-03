package handlers

import (
	"encoding/json"
	"medical-service/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetNotificationsByUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := mux.Vars(r)["userId"]
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	notifs, err := services.GetNotifications(uint(userId))
	if err != nil {
		http.Error(w, "Failed to get notifications", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notifs)
}

func MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	userIdStr := mux.Vars(r)["userId"]
	notifIdStr := mux.Vars(r)["notifId"]

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	notifId, err := strconv.Atoi(notifIdStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	if err := services.MarkNotificationAsRead(uint(userId), uint(notifId)); err != nil {
		http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
