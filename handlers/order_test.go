package handlers

import (
	"bytes"
	"encoding/json"
	"game-store-api/database"
	"game-store-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBuyProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := SetupRouter()

	product := models.Product{Name: "Zelda", Price: 6000, Stock: 1, SKU: "ZEL-1"}
	database.DB.Create(&product)

	user := models.User{Email: "gamer@test.com", Password: "hashed", Role: "user"}
	database.DB.Create(&user)

	token := GenerateTestToken(user.ID, "user")

	payload := map[string]uint{"product_id": product.ID}
	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var updatedProduct models.Product
	database.DB.First(&updatedProduct, product.ID)
	assert.Equal(t, 0, updatedProduct.Stock, "Stock should decrease by 1")

	var order models.Order
	err := database.DB.Where("user_id = ? AND product_id = ?", user.ID, product.ID).First(&order).Error
	assert.Nil(t, err, "Order record should be found")
}

func TestBuyOutOfStock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := SetupRouter()

	// Create a product with 0 stock
	product := models.Product{Name: "Roguelite Game", Price: 2549, Stock: 0, SKU: "ROGUE-1"}
	database.DB.Create(&product)

	// Dummy token with fake ID
	token := GenerateTestToken(1, "user")

	// Attempt to buy
	payload := map[string]uint{"product_id": product.ID}
	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Out of stock", response["error"])
}
