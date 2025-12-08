package handlers

import (
	"game-store-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	order, err := h.service.Checkout(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Order placed successfully",
		"order_id":   order.ID,
		"total_paid": order.TotalCents,
	})
}
