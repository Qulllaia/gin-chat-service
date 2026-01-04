package database

import (
	"fmt"

	"main/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sqlx.DB
}

func CreateConnection(config *config.Config) (*Database, error) {
	db, err := sqlx.Connect("postgres", connectionStringFormat(config))
	if err != nil {
		return nil, fmt.Errorf("Failed Connection %s", err.Error())
	}

	return &Database{DB: db}, nil
}

func connectionStringFormat(config *config.Config) string {
	if config.DB_PASSWORD != "" {
		return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_SSLMODE)
	}
	return fmt.Sprintf("user=%s dbname=%s sslmode=%s", config.DB_USER, config.DB_NAME, config.DB_SSLMODE)
}

