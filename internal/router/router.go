package router

import (
	"github.com/gin-gonic/gin"

	"nimo-backend/internal/auth"
	"nimo-backend/internal/handlers"
	"nimo-backend/internal/middleware"
)

func New(authHandler *auth.Handler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/google", authHandler.GoogleLogin)
	}

	protected := r.Group("/")
	protected.Use(middleware.RequireAuth(jwtSecret))
	{
		protected.GET("/me", authHandler.Me)
		protected.GET("/households", handlers.ListHouseholds)
		protected.GET("/households/:id/dashboard", handlers.GetHouseholdDashboard)
		protected.POST("/devices/claim", handlers.ClaimDevice)
	}

	return r
}
