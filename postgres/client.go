package postgres

import (
	"fmt"
	"log"
	"redis-learn/config"
	"redis-learn/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewClient(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	db.AutoMigrate(&models.Person{})

	fmt.Println("Connected to PostgreSQL")
	return db
}
