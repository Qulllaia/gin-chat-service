package dto

import "github.com/lib/pq"

type ChatListDTO struct {
	Name            *string       `db:"name"`
	ID              int           `db:"id"`
	Users           pq.Int64Array `db:"users"`
	Chat_type_id    string        `db:chat_type_id`
	Chat_background *string       `db:chat_background`
	User_id         *int64
	LastMessage     string
}

type ChatIDURI struct {
	ID int `uri:"id""`
}

type UserIDURI struct {
	ID int `uri:"id""`
}

type UsersIDList struct {
	IDs       []int64 `json: "ids"`
	GroupName string  `json: "GroupName"`
}

