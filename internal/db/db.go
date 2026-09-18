package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"nimo-backend/internal/models"
)

func Connect(path string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := database.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}

	return database, nil
}
