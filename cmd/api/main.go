package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "game-store-api/internal/grpc/payment"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"game-store-api/internal/handlers"
	"game-store-api/internal/middleware"
	"game-store-api/internal/models"
	"game-store-api/internal/repository"
	"game-store-api/internal/service"
	"game-store-api/internal/worker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load env variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on system environment variables")
	}

	// Connect to database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database connected successfully")

	// Run migrations
	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.CartItem{})
	if err != nil {
		slog.Error("Failed to migrate database", "error", err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		slog.Warn("Failed to connect to Redis", "error", err)
		redisClient = nil
	} else {
		slog.Info("Redis connected successfully")
	}

	// Start background worker
	if redisClient != nil {
		go worker.StartEmailWorker(redisClient)
	}

	// Connect to payment service
	paymentSvcAddr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if paymentSvcAddr == "" {
		paymentSvcAddr = "127.0.0.1:50051"
	}
	paymentConn, err := grpc.NewClient(paymentSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to payment service:", err)
	}
	defer paymentConn.Close()

	paymentClient := pb.NewPaymentServiceClient(paymentConn)

	// Dependency injection
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewCartRepository(db)

	authService := service.NewAuthService(userRepo, redisClient)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, cartRepo, paymentClient, db)

	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup router
	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)

		v1.GET("/products", productHandler.GetAllProducts)
		v1.GET("/products/:product_id", productHandler.GetProduct)

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/products", middleware.AdminOnly(), productHandler.CreateProduct)

			protected.GET("/cart", cartHandler.GetCart)
			protected.POST("/cart", cartHandler.AddToCart)
			protected.DELETE("/cart/:product_id", cartHandler.RemoveFromCart)

			protected.POST("/cart/checkout", orderHandler.Checkout)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to listen", "error", err)
			os.Exit(1)
		}
	}()
	slog.Info("Server is running", "port", 8080, "env", "development")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}
	slog.Info("Server exited properly")
}
