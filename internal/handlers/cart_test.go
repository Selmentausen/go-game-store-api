package handlers

import (
	"bytes"
	"encoding/json"
	"game-store-api/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddToCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deps := SetupTestDependencies()
	r := SetupRouter(deps)

	product := models.Product{Name: "Test", Price: 5000, Stock: 10, SKU: "TEST-1"}
	deps.DB.Create(&product)

	user := models.User{Email: "test@example.com", Password: "password123"}
	deps.DB.Create(&user)
	token := GenerateTestToken(user.ID, "user")

	payload := map[string]interface{}{
		"product_id": product.ID,
		"quantity":   1,
	}
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/cart", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var cartItem models.CartItem
	deps.DB.Where("user_id = ?", user.ID).First(&cartItem)
	assert.Equal(t, product.ID, cartItem.ProductID)
	assert.Equal(t, 1, cartItem.Quantity)
}
