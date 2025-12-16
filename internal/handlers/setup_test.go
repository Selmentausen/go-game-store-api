package handlers

import (
	"context"
	pb "game-store-api/internal/grpc/payment"
	"game-store-api/internal/middleware"
	"game-store-api/internal/models"
	"game-store-api/internal/repository"
	"game-store-api/internal/service"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type MockPaymentClient struct{}

func (m *MockPaymentClient) ProcessPayment(ctx context.Context, in *pb.PaymentRequest, opts ...grpc.CallOption) (*pb.PaymentResponse, error) {
	return &pb.PaymentResponse{
		Success:       true,
		TransactionId: "TEST_TXN_123",
		Message:       "Payment processed via Mock",
	}, nil
}

type TestDeps struct {
	DB             *gorm.DB
	AuthHandler    *AuthHandler
	ProductHandler *ProductHandler
	OrderHandler   *OrderHandler
	CartHandler    *CartHandler
}

func SetupTestDependencies() TestDeps {
	os.Setenv("JWT_SECRET", "test_secret_key")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to migrate test database: " + err.Error())
	}
	db.AutoMigrate(&models.Product{}, &models.User{}, &models.Order{}, &models.CartItem{}, &models.OrderItem{})

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewCartRepository(db)

	mockPayment := &MockPaymentClient{}

	authService := service.NewAuthService(userRepo, nil)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, cartRepo, mockPayment, db)

	return TestDeps{
		DB:             db,
		AuthHandler:    NewAuthHandler(authService),
		ProductHandler: NewProductHandler(productService),
		CartHandler:    NewCartHandler(cartService),
		OrderHandler:   NewOrderHandler(orderService),
	}
}

func SetupRouter(deps TestDeps) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", deps.AuthHandler.Register)
		v1.POST("/auth/login", deps.AuthHandler.Login)
		v1.GET("/products", deps.ProductHandler.GetProducts)

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/products", middleware.AdminOnly(), deps.ProductHandler.CreateProduct)

			protected.GET("/cart", deps.CartHandler.GetCart)
			protected.POST("/cart", deps.CartHandler.AddToCart)
			protected.DELETE("/cart/:product_id", deps.CartHandler.RemoveFromCart)

			protected.POST("/cart/checkout", deps.OrderHandler.Checkout)
		}
	}
	return r
}

func GenerateTestToken(userID uint, role string) string {
	claims := jwt.MapClaims{
		"sub":  float64(userID),
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test_secret_key"))
	return tokenString
}
