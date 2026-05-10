package database

import (
	"fmt"
	"log"

	"movie-night-planner-backend/internal/config"
	"movie-night-planner-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(cfg *config.DatabaseConfig) error {
	var err error

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db

	// Auto migrate schemas
	err = autoMigrate(db)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Evening{},
		&models.EveningFilm{},
		&models.Vote{},
		&models.Comment{},
	)
	if err != nil {
		return err
	}
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
