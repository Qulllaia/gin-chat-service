package models

type User struct {
	ID        int             `db:"id" json:"id"`
	Chat_list []int64 `db:"chat_list" json:"chat_list"`
}