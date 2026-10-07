package main

import (
	"log"

	"nimo-backend/internal/auth"
	"nimo-backend/internal/config"
	"nimo-backend/internal/db"
	"nimo-backend/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	database, err := db.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	authHandler := auth.NewHandler(database, cfg.JWTSecret, cfg.GoogleClientID)
	r := router.New(authHandler, cfg.JWTSecret, cfg.AllowedOrigins, cfg.AllowLocalhost)

	log.Printf("NIMO backend running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
