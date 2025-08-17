package queries

import (
	"main/database"
	. "main/database/models"
)

type AuthQuery struct {
	*database.Database
}

func AuthQueryConstructor(db *database.Database) *AuthQuery {
	return &AuthQuery{db}
}


func (aq *AuthQuery) GetUserByName(name string) (*User, error) {
    var user User
    err := aq.DB.QueryRow(`
        SELECT id, name, password 
        FROM "user" 
        WHERE name = $1
    `, name).Scan(&user.ID, &user.Name, &user.Password)
    
    if err != nil {
        return nil, err
    }
    return &user, nil
}