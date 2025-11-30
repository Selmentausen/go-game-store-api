package handlers

import (
	"game-store-api/database"
	"game-store-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	database.DB = db
	database.DB.AutoMigrate(&models.Product{}, &models.User{})
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/products", CreateProduct)
		v1.GET("/products", GetProduct)
		v1.POST("/auth/register", Register)
	}
	return r
}
