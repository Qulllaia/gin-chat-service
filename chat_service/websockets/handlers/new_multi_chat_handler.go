package handlers

import (
	"main/types"
	. "main/types"
	"main/websockets/handlers/base"

	"github.com/gorilla/websocket"
)

type NewMultiChatHandler struct {
	*base.BaseHandler	
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
	nmch.CreateChatContext(current_user_id, messageType, conn, 
		&map[string]interface{}{
			"type": "NEW_MULTIPLE_CHAT",
		}, a,
	)
	
	for _, value := range message.User_ids {
		if val, exists := userIdToConns[value]; exists {
			nmch.CreateChatContext(value, messageType, val.WS, 
				&map[string]interface{}{
					"type": "NEW_MULTIPLE_CHAT",
				}, a,
			)
		}
	}
}