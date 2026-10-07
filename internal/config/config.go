package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	JWTSecret      string
	GoogleClientID string
	DBPath         string
	AllowedOrigins []string
	AllowLocalhost bool
}

func Load() (Config, error) {
	_ = godotenv.Load()

	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must be set and at least 32 characters long")
	}

	var allowedOrigins []string
	for _, origin := range strings.Split(os.Getenv("CORS_ORIGINS"), ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowedOrigins = append(allowedOrigins, origin)
		}
	}

	return Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      jwtSecret,
		GoogleClientID: strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
		DBPath:         getEnv("DB_PATH", "nimo.db"),
		AllowedOrigins: allowedOrigins,
		AllowLocalhost: strings.EqualFold(os.Getenv("CORS_ALLOW_LOCALHOST"), "true"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
