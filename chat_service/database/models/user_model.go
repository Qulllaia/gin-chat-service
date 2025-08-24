package models

import "github.com/lib/pq"

type User struct {
	ID        int             `db:"id" json:"id"`
	Chat_list pq.Int64Array `db:"chat_list" json:"chat_list"`
}