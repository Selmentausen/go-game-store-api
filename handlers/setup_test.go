package handlers

import (
	"game-store-api/database"
	"game-store-api/middlewares"
	"game-store-api/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func SetupTestDB() {
	os.Setenv("JWT_SECRET", "test_secret_key")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	database.DB = db
	err = database.DB.AutoMigrate(&models.Product{}, &models.User{}, &models.Order{})
	if err != nil {
		panic("Failed to migrate test database: " + err.Error())
	}
}

func GenerateTestToken(userID uint, role string) string {
	claims := jwt.MapClaims{
		"sub":  float64(userID),
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test_secret_key"))
	return tokenString
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", Login)
		v1.POST("/auth/register", Register)
		v1.GET("/products", GetProduct)

		protected := v1.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("/products", middlewares.AdminOnly(), CreateProduct)
			protected.POST("/orders", CreateOrder)
		}
	}
	return r
}
