package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deps := SetupTestDependencies()
	r := SetupRouter(deps)

	token := GenerateTestToken(1, "admin")

	payload := map[string]interface{}{
		"name":        "Test Game",
		"description": "Desc",
		"price":       1000,
		"stock":       10,
		"sku":         "TEST-001",
	}

	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
