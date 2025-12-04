package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game-store-api/internal/models"
	"game-store-api/internal/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
}

func NewAuthService(userRepo repository.UserRepository, redisClient *redis.Client) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

func (s *AuthService) Register(email, password string) error {
	existing, _ := s.userRepo.GetUserByEmail(email)
	if existing.ID != 0 {
		return errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    email,
		Password: string(hashed),
		Role:     "user",
	}

	if err := s.userRepo.CreateUser(&user); err != nil {
		return err
	}

	if s.redisClient != nil {
		taskPayload := map[string]string{
			"email":   user.Email,
			"user_id": fmt.Sprintf("%d", user.ID),
			"type":    "welcome_email",
		}
		jsonBody, _ := json.Marshal(taskPayload)
		s.redisClient.RPush(context.Background(), "send_email_queue", jsonBody)
	}
	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"sub":  float64(user.ID),
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}
