package services

import (
	"school-service/database"
	"school-service/models"
)

func AddNotification(notif models.Notification) error {
	return database.DB.Create(&notif).Error
}

func GetNotificationsByStudent(studentID uint) ([]models.Notification, error) {
	var notifs []models.Notification
	if err := database.DB.Where("user_id = ?", studentID).Find(&notifs).Error; err != nil {
		return nil, err
	}
	return notifs, nil
}

func MarkNotificationAsRead(userID uint, notifID uint) error {
	var notif models.Notification
	if err := database.DB.First(&notif, notifID).Error; err != nil {
		return err
	}

	if notif.UserId != userID {
		return nil
	}

	notif.Read = true
	return database.DB.Save(&notif).Error
}
