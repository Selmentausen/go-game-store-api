package main

import (
	"fmt"
	"game-store-api/database"
	"game-store-api/models"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase()

	adminEmail := "admin@gamestore.com"
	adminPassword := "admin123"

	var existingUser models.User
	if err := database.DB.Where("email = ?", adminEmail).First(&existingUser).Error; err == nil {
		fmt.Println("Admin user already exists")
		return
	}

	admin := models.User{
		Email:    adminEmail,
		Password: adminPassword,
		Role:     "admin",
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to created admin: %v", err)
	}

	fmt.Println("Admin created successfully!")
	fmt.Println("Email: %s\nPassword: %s\n", adminEmail, adminPassword)
}
