package queries

import "main/database"

type UserQuery struct {
	*database.Database
}

func UserQueryConstructor(db *database.Database) *UserQuery {
	return &UserQuery{db}
}

func (uq *UserQuery) InsertUser(name string, password string) error {
	_, err := uq.DB.Exec(`INSERT INTO "user" (name, password) VALUES ($1, $2)`, name, password);
	return err; 
}