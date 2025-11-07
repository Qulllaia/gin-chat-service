package queries

import (
	"database/sql"
	"main/database"
	. "main/database/models"
)

type AuthQuery struct {
	*database.Database
}

func AuthQueryConstructor(db *database.Database) *AuthQuery {
	return &AuthQuery{db}
}

func (aq *AuthQuery) GetUserByNameOrEmail(name, email string) (*User, error) {
	var user User
	err := aq.DB.QueryRow(`
        SELECT id, email, name, password 
        FROM "user" 
        WHERE name = $1 or email = $2
    `, name, email).Scan(&user.ID, &user.Email, &user.Name, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
