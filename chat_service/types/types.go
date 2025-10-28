package types

import (
	"main/database/queries"
	"main/redis"

	"github.com/gin-gonic/gin"
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

type Endpoint[T any] func(context *gin.Context) (ErrorType, T, error)

type HttpResponse struct {
	Done   bool        `json:"done"`
	Result interface{} `json:"result,omitempty"`
	ErrorType ErrorType `json:"error_type,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type ImageResponse struct {
	Message string `json:"message"`
	Filename string `json:"filename"`
	Url string `json:"url"`
	Full_url string `json:"full_url"`
}