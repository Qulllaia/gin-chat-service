package handlers

import (
	"encoding/json"
	"main/types"
	. "main/types"

	"github.com/gorilla/websocket"
)

type NewMultiChatHandler struct {
}

func NewNewMultiChatHandler() types.Handler {
	return &NewMultiChatHandler{}
}

func (nmch *NewMultiChatHandler) Handle(message MessageWS, messageType int, conn *websocket.Conn, actor Actor) {
	nmch.broadcastChatCreationNotify(message, messageType, conn, actor)
}

func (nmch *NewMultiChatHandler) broadcastChatCreationNotify(message MessageWS, messageType int, conn *websocket.Conn, a Actor) {
	connsToUserId := a.GetConnectoinsToUsers()
	current_user_id := connsToUserId[conn];
	userIdToConns := a.GetUserConnections()
	nmch.createChatContext(current_user_id, messageType, conn, 
		&map[string]interface{}{
			"type": "NEW_MULTIPLE_CHAT",
		}, a,
	)
	
	for _, value := range message.User_ids {
		if val, exists := userIdToConns[value]; exists {
			nmch.createChatContext(value, messageType, val.WS, 
				&map[string]interface{}{
					"type": "NEW_MULTIPLE_CHAT",
				}, a,
			)
		}
	}
}

func (nmch *NewMultiChatHandler) createChatContext(user_id int, messageType int, conn *websocket.Conn, response_data *map[string]interface{}, a Actor) {

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