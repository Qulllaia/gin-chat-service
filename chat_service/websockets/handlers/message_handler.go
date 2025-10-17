package handlers

import (
	"encoding/json"
	"main/types"
	"main/utils"
	"strconv"

	"github.com/gorilla/websocket"
)

type MessageHandler struct {
}

func NewMessageHandler() types.Handler {
	return &MessageHandler{}
}

func (mh *MessageHandler) Handle(message types.MessageWS, messageType int, conn *websocket.Conn, actor types.Actor) {
	mh.broadcastMessage(message, messageType, conn, actor)
}

func (mh *MessageHandler) broadcastMessage(message types.MessageWS, messageType int, conn *websocket.Conn, actor types.Actor) {
	var getter_conn *websocket.Conn
	user_conns := actor.GetUserConnections()
	wsq := actor.GetWSQ()
	for user_id, value := range user_conns {
		chat_id, _ := strconv.Atoi(message.Chat_id)

		if utils.Contains(value.Chat_ids, int64(chat_id)) {
			getter_conn = value.WS
			if getter_conn != conn {
				responseData, err := json.Marshal(map[string]interface{}{
					"message": message.Message,
					"chat_id": chat_id,
					"type":    "MESSAGE",
				})
				if err = getter_conn.WriteMessage(messageType, []byte(responseData)); err != nil {
					println(err.Error())
				}
			} else {
				if err := wsq.InsertMessageIntoChatHistory(chat_id, user_id, string(message.Message)); err != nil {
					println(err.Error())
				}
			}
		}
	}
}
