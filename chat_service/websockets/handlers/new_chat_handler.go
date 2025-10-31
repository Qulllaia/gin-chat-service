package handlers

import (
	"main/types"
	. "main/types"
	"main/websockets/handlers/base"
	"strconv"

	"github.com/gorilla/websocket"
)

type NewChatHandler struct {
	*base.BaseHandler
}

func NewNewChatHandler() types.Handler {
	return &NewChatHandler{}
}

func (nch *NewChatHandler) Handle(message MessageWS, messageType int, conn *websocket.Conn, actor Actor) {
	nch.createChatWithMessage(message, conn, messageType, actor);
	nch.BroadcastMessage(message, messageType, conn, actor);
}

func (nch *NewChatHandler) createChatWithMessage(msg MessageWS, conn *websocket.Conn, messageType int, a Actor) {
	actorRedisConnector := a.GetRDB();	
	actorWSQ := a.GetWSQ()	
	userConns := a.GetUserConnections()	
	connsToUserIds := a.GetConnectoinsToUsers()
	
	current_user_id := connsToUserIds[conn];
	to_user_id, _ := strconv.Atoi(msg.User_id);

	err := actorRedisConnector.DeleteData(strconv.Itoa(current_user_id))
	
	err = actorRedisConnector.DeleteData(msg.User_id)

	if err != nil {
		println("Error while deleting redis", err.Error());
	}
	

	err, chat_id := actorWSQ.CreateChatWithMessage(current_user_id, to_user_id);
	if err != nil {
		println("Error while creating a chat", err.Error());
	}
 
	if err != nil {
		println("Ошибка формирования ответа при создании нового чата", err.Error())
	}

	actorWSQ.InsertMessageIntoChatHistory(chat_id, current_user_id, msg.Message);

	nch.CreateChatContext(current_user_id, messageType, conn, &map[string]interface{} {
		"type": NEW_CHAT,
		"chat_id": chat_id,
	},a)

	if toUserData, exists := userConns[to_user_id]; exists  {
		
		nch.CreateChatContext(to_user_id, messageType, toUserData.WS, &map[string]interface{} {
			"type": NEW_CHAT,
			"chat_id": chat_id,

		},a)
	}
}