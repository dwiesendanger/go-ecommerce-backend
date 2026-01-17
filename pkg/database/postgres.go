package database

import (
	"fmt"
	"log"
	"os"

	"ecommerce-platform/internal/core/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "user"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "ecommerce"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Running AutoMigrate...")
	err = db.AutoMigrate(&domain.User{}, &domain.Product{}, &domain.Cart{}, &domain.CartItem{}, &domain.Order{}, &domain.OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}
