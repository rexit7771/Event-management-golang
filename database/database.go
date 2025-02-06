package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Connect() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		// *Env untuk Deploy railway
		// os.Getenv("PGHOST"),
		// os.Getenv("PGPORT"),
		// os.Getenv("PGUSER"),
		// os.Getenv("PGPASSWORD"),
		// os.Getenv("PGDATABASE"),

		// * Env untuk local
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}
	fmt.Println("Database Connected")
}
