package models

type Message struct {
	Id               int64  `db:"id" json:"id"`
	Message          string `db:"message" json:"message"`
	User_id          int64  `db:"user_id" json:"user_id"`
	Chat_id          int64  `db:"chat_id" json:"chat_id"`
	Timestamp        string `db:"timestamp" json:"timestamp"`
	IsThisUserSender bool
}