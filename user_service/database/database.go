package database

import (
	"context"
	"fmt"
	"main/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sqlx.DB
}

func CreateConnection(context context.Context) (*Database, error) {
	config, exists := context.Value("config").(*config.Config);
	if !exists {
		return nil, fmt.Errorf("Config data has not been setted")
	}

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_SSLMODE)

	db, err := sqlx.Connect("postgres", connectionString);
	if err != nil{
		println("Failed Connection", err)
	}
	return &Database{DB: db}, nil;
}