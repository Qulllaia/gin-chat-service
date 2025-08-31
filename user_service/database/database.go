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

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_SSLMODE)

	db, err := sqlx.Connect("postgres", connectionString);
	if err != nil{
		println("Failed Connection", err)
	}
	return &Database{DB: db}, nil;
}