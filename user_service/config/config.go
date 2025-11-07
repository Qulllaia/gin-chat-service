package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_USER        string
	DB_PASSWORD    string
	DB_NAME        string
	DB_SSLMODE     string
	JWT_SECRET     string
	EMAIL_SENDER   string
	EMAIL_PASSWORD string
	SMTP_ADDR      string
	SMTP_PORT      string
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
	emailSender := os.Getenv("EMAIL_SENDER")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	smtpAddr := os.Getenv("SMTP_ADDR")
	smtpPort := os.Getenv("SMTP_PORT")

	return &Config{
		DB_USER:        dbUser,
		DB_PASSWORD:    dbPassword,
		DB_NAME:        dbName,
		DB_SSLMODE:     dbSSLMode,
		JWT_SECRET:     jwtSecret,
		EMAIL_SENDER:   emailSender,
		EMAIL_PASSWORD: emailPassword,
		SMTP_ADDR:      smtpAddr,
		SMTP_PORT:      smtpPort,
	}, nil
}
