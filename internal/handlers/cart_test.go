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

func TestCartLifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deps := SetupTestDependencies()
	r := SetupRouter(deps)

	product := models.Product{Name: "Test", Price: 5000, Stock: 10, SKU: "TEST-1"}
	deps.DB.Create(&product)

	user := models.User{Email: "test@example.com", Password: "password123"}
	deps.DB.Create(&user)
	token := GenerateTestToken(user.ID, "user")

	sendCartRequest := func(qty int) *httptest.ResponseRecorder {
		payload := map[string]interface{}{"product_id": product.ID, "quantity": qty}
		jsonVal, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/cart", bytes.NewBuffer(jsonVal))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	// TEST 1: Add product
	w1 := sendCartRequest(1)
	assert.Equal(t, http.StatusOK, w1.Code)

	var item models.CartItem
	deps.DB.Where("user_id = ?", user.ID).First(&item)
	assert.Equal(t, 1, item.Quantity, "Should start with quantity 1")

	// TEST 2: Increment product
	w2 := sendCartRequest(1)
	assert.Equal(t, http.StatusOK, w2.Code)
	deps.DB.Where("user_id = ?", user.ID).First(&item)
	assert.Equal(t, 2, item.Quantity, "1 + 1 should be 2")

	// TEST 3: Decrement product
	w3 := sendCartRequest(-1)
	assert.Equal(t, http.StatusOK, w3.Code)
	deps.DB.Where("user_id = ?", user.ID).First(&item)
	assert.Equal(t, 1, item.Quantity, "2 - 1 should be 1")

	// TEST 4: Remove product via decrement
	w4 := sendCartRequest(-1)
	assert.Equal(t, http.StatusOK, w4.Code)
	err := deps.DB.Where("user_id = ?", user.ID).First(&item).Error
	assert.Error(t, err, "Record should be deleted when quantity hits 0")
}

func TestNegativeInitialAdd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deps := SetupTestDependencies()
	r := SetupRouter(deps)

	product := models.Product{Name: "Test", Price: 5000, Stock: 10, SKU: "TEST-1"}
	deps.DB.Create(&product)

	user := models.User{Email: "test@example.com", Password: "password123"}
	deps.DB.Create(&user)
	token := GenerateTestToken(user.ID, "user")

	payload := map[string]interface{}{"product_id": product.ID, "quantity": -1}
	jsonVal, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/cart", bytes.NewBuffer(jsonVal))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var count int64
	deps.DB.Model(&models.CartItem{}).Count(&count)
	assert.Equal(t, int64(0), count, "Should not create cart item with negative quantity")
}

func TestRemoveItemEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deps := SetupTestDependencies()
	r := SetupRouter(deps)

	product := models.Product{Name: "Test", Price: 5000, Stock: 10, SKU: "TEST-1"}
	deps.DB.Create(&product)
	user := models.User{Email: "test@example.com", Password: "password123"}
	deps.DB.Create(&user)
	token := GenerateTestToken(user.ID, "user")
	deps.DB.Create(&models.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 5})

	req, _ := http.NewRequest("DELETE", "/api/v1/cart/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var count int64
	deps.DB.Model(&models.CartItem{}).Count(&count)
	assert.Equal(t, int64(0), count)

}
