package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"game-store-api/database"
	"game-store-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := SetupRouter()

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "mysecretpassword",
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Check if user is created
	assert.Equal(t, http.StatusCreated, w.Code)

	// Check if password is protected and not stored in plain text
	var user models.User
	database.DB.Where("email = ?", "test@example.com").First(&user)

	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEqual(t, "mysecretpassword", user.Password, "Password should be hashed!")
	assert.Len(t, user.Password, 60, "Bcrypt hash should be 60 chars long")
}

func TestRegisterDuplicateEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := SetupRouter()

	database.DB.Create(&models.User{
		Email:    "duplicate@example.com",
		Password: "mysecretpassword",
	})

	payload := map[string]string{
		"email":    "duplicate@example.com",
		"password": "newpassword",
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
