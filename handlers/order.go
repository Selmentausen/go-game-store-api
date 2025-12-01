package handlers

import (
	"game-store-api/database"
	"game-store-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type OrderInput struct {
	ProductID uint `json:"product_id" binding:"required"`
}

func CreateOrder(c *gin.Context) {
	var input OrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("UserID")
	tx := database.DB.Begin()

	// check Product existende and lock the row (preven race conditions)
	var product models.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, input.ProductID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.Stock <= 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Out of stock"})
		return
	}

	order := models.Order{
		UserID:    userID.(uint),
		ProductID: input.ProductID,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	product.Stock = product.Stock - 1
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "Order placed successfully", "order_id": order.ID})
}
