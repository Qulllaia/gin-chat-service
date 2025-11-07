package dto

type UserDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	ID       int    `json:"id"`
}

type UserIDURI struct {
	ID int `uri:"id""`
}

type UserWithoutPasswordDTO struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type TokenUriDto struct {
	Token string `uri:"token"`
}
