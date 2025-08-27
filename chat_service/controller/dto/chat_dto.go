package dto

type ChatListDTO struct {
	Name string `db:"name"`
	ID   int    `db:"id"`
}

type ChatIDURI struct {
	ID int `uri:"id""`
}

type UserIDURI struct {
	ID int `uri:"id""`
}