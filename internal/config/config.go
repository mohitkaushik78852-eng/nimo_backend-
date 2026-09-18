package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	JWTSecret      string
	GoogleClientID string
	DBPath         string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-me"),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		DBPath:         getEnv("DB_PATH", "nimo.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
