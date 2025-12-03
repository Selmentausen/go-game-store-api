package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminRouteProtection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := SetupRouter()

	productPayload := []byte(`{"name":"Hacked Game", "price":100, "stock":10, "sku":"HACK-1"}`)

	// Normal User tries to add product
	userToken := GenerateTestToken(1, "user")
	req1, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(productPayload))
	req1.Header.Set("Authorization", "Bearer "+userToken)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// Should be foridden (403)
	assert.Equal(t, http.StatusForbidden, w1.Code)

	// Admin tries to add product
	adminToken := GenerateTestToken(2, "admin")

	req2, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(productPayload))
	req2.Header.Set("Authorization", "Bearer "+adminToken)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusCreated, w2.Code)
}
