package main

import (
	"fmt"
	"game-store-api/internal/models"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on system environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	adminEmail := "admin@gamestore.com"
	adminPassword := "admin123"

	var existingUser models.User
	if err := db.Where("email = ?", adminEmail).First(&existingUser).Error; err == nil {
		slog.Info("Admin user already exists", "email", adminEmail)
		return
	}

	admin := models.User{
		Email:    adminEmail,
		Password: adminPassword,
		Role:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		slog.Error("Failed to created admin", "error", err)
		os.Exit(1)
	}

	slog.Info("Admin created successfully!")
}
