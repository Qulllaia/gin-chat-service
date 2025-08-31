package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_USER string;
	DB_PASSWORD string;
	DB_NAME string;
	DB_SSLMODE string;
	JWT_SECRET string;
}

func CreateConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
        println("Error loading .env file:", err)
		return nil, err
    }

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	jwtSecret := os.Getenv("JWT_SECRET")


	return  &Config{
		DB_USER: dbUser,
		DB_PASSWORD: dbPassword,
		DB_NAME: dbName,
		DB_SSLMODE: dbSSLMode,
		JWT_SECRET: jwtSecret,
	}, nil;
}