package dto

type UserDTO struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	ID       int    `json:"id"`
}

type UserIDURI struct {
	ID int `uri:"id""`
}