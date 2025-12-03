package main

import (
	"context"
	"fmt"
	"game-store-api/database"
	"game-store-api/models"
	"game-store-api/routes"
	"game-store-api/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase()
	database.ConnectRedis()
	database.DB.AutoMigrate(&models.Product{}, &models.User{}, &models.Order{})

	go worker.StartEmailWorker()

	r := gin.Default()
	r.Static("/static", "./static")
	routes.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Run Server in a separate Goroutine so it doesn't block the main thread.
	// This allows the main thread to listen for termination signals.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	fmt.Println("Server is running on port 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("\nShutting down server...")

	// When a kill signal is received, we give the server 5 seconds to finish
	// processing currently active requests before forcibly closing.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exited properly")
}
