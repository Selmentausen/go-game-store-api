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

func TestCheckoutFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deps := SetupTestDependencies()
	r := SetupRouter(deps)

	product := models.Product{Name: "Zelda", Price: 6000, Stock: 10, SKU: "ZEL-1"}
	deps.DB.Create(&product)

	user := models.User{Email: "gamer@test.com", Password: "hashed", Role: "user"}
	deps.DB.Create(&user)

	token := GenerateTestToken(user.ID, "user")

	// Test adding products to cart
	cartPayload := map[string]interface{}{
		"product_id": product.ID, "quantity": 2,
	}
	jsonBody, _ := json.Marshal(cartPayload)

	req1, _ := http.NewRequest("POST", "/api/v1/cart", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)

	// Test checkout of cart
	req2, _ := http.NewRequest("POST", "/api/v1/cart/checkout", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusCreated, w2.Code)

	// Check order created
	var order models.Order
	err := deps.DB.Preload("Items").Where("user_id = ?", user.ID).First(&order).Error
	assert.Nil(t, err)
	assert.Equal(t, 12000, order.TotalCents)
	assert.Equal(t, 1, len(order.Items))

	// Check stock deducted from product
	var updatedProduct models.Product
	deps.DB.First(&updatedProduct, product.ID)
	assert.Equal(t, 8, updatedProduct.Stock, "Stock should decrease by 1")

	var cartItems []models.CartItem
	deps.DB.Where("user_id = ?", user.ID).Find(&cartItems)
	assert.Equal(t, 0, len(cartItems))
}

func TestCheckoutEmptyCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deps := SetupTestDependencies()
	r := SetupRouter(deps)

	user := models.User{Email: "test@example.com", Password: "hashed", Role: "user"}
	deps.DB.Create(&user)
	token := GenerateTestToken(user.ID, "user")

	// Try checkout with empty cart
	req, _ := http.NewRequest("POST", "/api/v1/cart/checkout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
