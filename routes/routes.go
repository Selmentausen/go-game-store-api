package routes

import (
	"game-store-api/handlers"
	"game-store-api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", handlers.Register)
		v1.POST("/auth/login", handlers.Login)
		v1.GET("/products", handlers.GetProduct)

		protected := v1.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("/products", middlewares.AdminOnly(), handlers.CreateProduct)
		}
	}
}
