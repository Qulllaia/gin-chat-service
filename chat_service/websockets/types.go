package websockets

import (
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

type MessageType string

type MessageWS struct {
	Chat_id string `json:"chat_id"`
	Message string `json:"messages"`
}

type UserWSData struct {
	chat_ids pq.Int64Array
	ws       *websocket.Conn
}

type mailboxObject struct {
	message MessageWS
	conn *websocket.Conn
	MessageType int
}