package database

import (
	"log"
	"os"

	"auth-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database: ", err)
	}

	err = DB.AutoMigrate(&models.Member{})
	if err != nil {
		log.Fatal("❌ Failed to migrate database: ", err)
	}

	log.Println("✅ Database connected and migrated")
}
