package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func CreateConfigDatabase() string {
	if err := godotenv.Load(); err != nil {
        println("Error loading .env file:", err)
    }

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUser, dbPassword, dbName, dbSSLMode)

	return  connStr;
}