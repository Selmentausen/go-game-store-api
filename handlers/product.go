package handlers

import (
	"game-store-api/database"
	"game-store-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var input models.Product

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&input)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func GetProduct(c *gin.Context) {
	var products []models.Product

	database.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}
