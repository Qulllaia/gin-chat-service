package dto

import "github.com/lib/pq"

type ChatListDTO struct {
	Name      *string `db:"name"`
	ID        int     `db:"id"`
	Users     pq.Int64Array `db:"users"`
	Chat_type string  `db:chat_type`
}

type ChatIDURI struct {
	ID int `uri:"id""`
}

type UserIDURI struct {
	ID int `uri:"id""`
}

type UsersIDList struct {
	IDs []int64 `json: "ids"`
	GroupName string `json: "GroupName"`
}