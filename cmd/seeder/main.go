package main

import (
	"fmt"
	"log/slog"
	"os"

	"game-store-api/internal/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("No .env file found")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting Database Seed...")

	slog.Info("Cleaning old data...")
	db.Exec("DELETE FROM order_items")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM cart_items")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []models.User{
		{Email: "admin@gamestore.com", Password: string(hashedPass), Role: "admin"},
		{Email: "player1@test.com", Password: string(hashedPass), Role: "user"},
		{Email: "player2@test.com", Password: string(hashedPass), Role: "user"},
	}

	if err := db.Create(&users).Error; err != nil {
		slog.Error("Failed to create users", "error", err)
		os.Exit(1)
	}
	slog.Info("Users seeded (password: password123)")

	products := []models.Product{
		{

			Name:        "Elden Ring",
			Description: "THE NEW FANTASY ACTION RPG. Rise, Tarnished, and be guided by grace to brandish the power of the Elden Ring and become an Elden Lord in the Lands Between.",
			Price:       5999, // $59.99
			Stock:       50,
			SKU:         "ELD-001",
		},
		{
			Name:        "Hollow Knight",
			Description: "Forge your own path in Hollow Knight! An epic action adventure through a vast ruined kingdom of insects and heroes.",
			Price:       1499, // $14.99
			Stock:       100,
			SKU:         "HK-002",
		},
		{
			Name:        "Cyberpunk 2077",
			Description: "Cyberpunk 2077 is an open-world, action-adventure RPG set in the dark future of Night City — a dangerous megalopolis obsessed with power, glamor, and ceaseless body modification.",
			Price:       2999,
			Stock:       25,
			SKU:         "CP-2077",
		},
		{
			Name:        "Stardew Valley",
			Description: "You've inherited your grandfather's old farm plot in Stardew Valley. Armed with hand-me-down tools and a few coins, you set out to begin your new life.",
			Price:       1499,
			Stock:       200,
			SKU:         "SDV-004",
		},
		{
			Name:        "Satisfactory",
			Description: "For the avid enjoyers of spaghetti and massive timid creatures. Try your best to build a factory that is somewhat comprehensible",
			Price:       6999,
			Stock:       10, // Low stock to test "Sold Out"
			SKU:         "GOW-RAG",
		},
		{
			Name:        "Minecraft",
			Description: "Prepare for an adventure of limitless possibilities as you build, mine, battle mobs, and explore the ever-changing Minecraft landscape.",
			Price:       2999,
			Stock:       500,
			SKU:         "MC-001",
		},
		{
			Name:        "Half-Life 3",
			Description: "The game that never came out. Extremely rare.",
			Price:       99999,
			Stock:       0, // Out of stock test
			SKU:         "HL3-CONFIRMED",
		},
	}

	if err := db.Create(&products).Error; err != nil {
		slog.Error("Failed to create products", "error", err)
		os.Exit(1)
	}
	slog.Info("Products seeded")

	slog.Info("Seeding Complete!")
}
