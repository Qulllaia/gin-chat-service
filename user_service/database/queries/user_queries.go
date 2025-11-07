package queries

import (
	"main/controller/dto"
	"main/database"
	. "main/database/models"

	"github.com/lib/pq"
)

type UserQuery struct {
	*database.Database
}

func UserQueryConstructor(db *database.Database) *UserQuery {
	return &UserQuery{db}
}

func (uq *UserQuery) InsertUser(email, name, password string) (int64, error) {
	var id int64
	err := uq.DB.QueryRow(`INSERT INTO "user" (email, name, password) VALUES ($1, $2, $3) RETURNING id`, email, name, password).Scan(&id)
	return id, err
}

func (uq *UserQuery) GetUserByID(id int) (*User, error) {
	var user User
	err := uq.DB.QueryRow(`
        SELECT id, name, password 
        FROM "user" 
        WHERE id = $1
    `, id).Scan(&user.ID, &user.Name, &user.Password)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (uq *UserQuery) GetUserNamesByIDs(ids []int) (map[int]string, error) {

	rows, err := uq.DB.Query(`
        SELECT id, name 
        FROM "user" 
        WHERE id = ANY($1)
    `, pq.Array(ids))

	userIdUserName := make(map[int]string)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		userIdUserName[user.ID] = user.Name
	}

	return userIdUserName, nil
}

func (uq *UserQuery) GetAllUsers() ([]User, error) {
	rows, err := uq.DB.Query(`
        SELECT id, name, password 
        FROM "user"
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (uq *UserQuery) UpdateUser(id int, name, password string) error {
	_, err := uq.DB.Exec(`
        UPDATE "user" 
        SET name = $1, password = $2 
        WHERE id = $3
    `, name, password, id)
	return err
}

func (uq *UserQuery) DeleteUser(id int) error {
	_, err := uq.DB.Exec(`
        DELETE FROM "user" 
        WHERE id = $1
    `, id)
	return err
}

func (uq *UserQuery) GetUserExceptCurrent(id int) (*[]dto.UserWithoutPasswordDTO, error) {
	rows, err := uq.DB.Query(`
        SELECT id, name
        FROM "user" 
        WHERE id != $1
    `, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []dto.UserWithoutPasswordDTO
	for rows.Next() {
		var user dto.UserWithoutPasswordDTO
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}
