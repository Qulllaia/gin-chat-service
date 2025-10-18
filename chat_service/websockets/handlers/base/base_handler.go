package base

import (
	"encoding/json"
	"main/types"
	"main/utils"
	"strconv"

	"github.com/gorilla/websocket"
)

type BaseHandler struct {
}

func (bh *BaseHandler) BroadcastMessage(message types.MessageWS, messageType int, conn *websocket.Conn, actor types.Actor) {
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

func (bh *BaseHandler) CreateChatContext(user_id int, messageType int, conn *websocket.Conn, response_data *map[string]interface{}, a types.Actor) {

	a.AddClient(conn, user_id)

	if response_data != nil {
		responseData, err := json.Marshal(response_data)

		if err != nil {
			println(err.Error())
		}

		if err = conn.WriteMessage(messageType, responseData); err != nil {
			println(err.Error());	
		}
	}
}