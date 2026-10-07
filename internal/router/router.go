package router

import (
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"nimo-backend/internal/auth"
	"nimo-backend/internal/handlers"
	"nimo-backend/internal/middleware"
)

func New(authHandler *auth.Handler, jwtSecret string, allowedOrigins []string, allowLocalhost bool) *gin.Engine {
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		origin := c.GetHeader("Origin")
		if origin != "" && isAllowedOrigin(origin, allowedOrigins, allowLocalhost) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Vary", "Origin")
		}
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
	authGroup.Use(middleware.RateLimit(10, time.Minute))
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

func isAllowedOrigin(origin string, allowedOrigins []string, allowLocalhost bool) bool {
	if slices.Contains(allowedOrigins, strings.TrimSpace(origin)) {
		return true
	}
	if !allowLocalhost {
		return false
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme != "http" || parsed.Port() == "" {
		return false
	}
	return parsed.Hostname() == "localhost" || parsed.Hostname() == "127.0.0.1" || parsed.Hostname() == "[::1]"
}
