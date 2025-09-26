package database

import (
	"log"
	"os"
	"school-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database: ", err)
	}

	err = DB.AutoMigrate(
		&models.Absence{},
		&models.Class{},
		&models.Grade{},
		&models.Parent{},
		&models.Student{},
		&models.Subject{},
		&models.Teacher{},
	)
	if err != nil {
		log.Fatal("❌ Failed to migrate database: ", err)
	}

	log.Println("✅ School database connected and migrated")
	return DB
}
