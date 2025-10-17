package types

import (
	"main/database/queries"
	"main/redis"

	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

type MessageType string

type MessageWS struct {
	Chat_id  string `json:"chat_id"`
	User_id  string `json:"user_id"`
	User_ids []int  `json:"user_ids"`
	Message  string `json:"messages"`
	Type     string `json:"type"`
}

type UserWSData struct {
	Chat_ids pq.Int64Array
	WS       *websocket.Conn
}

type MailboxObject struct {
	Message     MessageWS
	Conn        *websocket.Conn
	MessageType int
}

type Actor interface {
	GetUserConnections() map[int]UserWSData
	GetWSQ() *queries.WSQueries
	GetRDB() *redis.RedisConnector
	GetConnectoinsToUsers() map[*websocket.Conn]int
	AddClient(conn *websocket.Conn, user_id int)
}

type Handler interface {
	Handle(message MessageWS, messageType int, conn *websocket.Conn, actor Actor)
}