package services

import (
	"medical-service/database"
	"medical-service/models"
)

func CreateNotification(userId uint, message string) error {
	notif := models.Notification{
		UserId:  userId,
		Message: message,
	}
	return database.DB.Create(&notif).Error
}

func GetNotifications(userId uint) ([]models.Notification, error) {
	var notifs []models.Notification
	err := database.DB.Where("user_id = ?", userId).Order("read ASC, created_at DESC").Find(&notifs).Error
	return notifs, err
}

func MarkNotificationAsRead(userId uint, notifId uint) error {
	return database.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notifId, userId).
		Update("read", true).Error
}
