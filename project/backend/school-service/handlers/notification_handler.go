package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"school-service/models"
	"school-service/services"
	"strconv"

	"github.com/gorilla/mux"
)

type NotificationDTO struct {
	StudentID uint   `json:"studentId"`
	Message   string `json:"message"`
}

func ReceiveNotification(w http.ResponseWriter, r *http.Request) {
	var notif NotificationDTO
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	log.Printf("[Notification] Student %d received message: %s", notif.StudentID, notif.Message)

	err := services.AddNotification(models.Notification{
		UserId:  notif.StudentID,
		Message: notif.Message,
	})
	if err != nil {
		http.Error(w, "Failed to save notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetNotificationsForStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentIDStr := vars["studentID"]
	studentID, err := strconv.ParseUint(studentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid studentID", http.StatusBadRequest)
		return
	}

	notifications, err := services.GetNotificationsByStudent(uint(studentID))
	if err != nil {
		http.Error(w, "Failed to get notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func MarkNotificationReadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["userID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	notifID, err := strconv.ParseUint(vars["notifID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	if err := services.MarkNotificationAsRead(uint(userID), uint(notifID)); err != nil {
		http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
