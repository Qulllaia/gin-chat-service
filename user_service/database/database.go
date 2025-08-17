package database

import (
	"main/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sqlx.DB
}

func CreateConnection() *Database {
	connectionString := config.CreateConfigDatabase();
	db, err := sqlx.Connect("postgres", connectionString);
	if err != nil{
		println("Failed Connection", err)
	}
	return &Database{DB: db};
}