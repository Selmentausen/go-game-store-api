package handlers

import (
	"bytes"
	"encoding/json"
	"game-store-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := SetupRouter()

	product := models.Product{
		Name:  "Test Game",
		Price: 1000,
		SKU:   "TEST-001",
	}

	adminToken := GenerateTestToken(1, "admin")
	jsonValue, _ := json.Marshal(product)

	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
