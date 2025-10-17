package handlers

import (
	"encoding/json"
	"main/types"
	. "main/types"
	"main/utils"
	"strconv"

	"github.com/gorilla/websocket"
)

type NewChatHandler struct {
}

func NewNewChatHandler() types.Handler {
	return &NewChatHandler{}
}

func (nch *NewChatHandler) Handle(message MessageWS, messageType int, conn *websocket.Conn, actor Actor) {
	nch.createChatWithMessage(message, conn, messageType, actor);
	nch.broadcastMessage(message, messageType, conn, actor);
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

	nch.createChatContext(current_user_id, messageType, conn, &map[string]interface{} {
		"type": "NEW_CHAT",
		"chat_id": chat_id,
	},a)

	if toUserData, exists := userConns[to_user_id]; exists  {
		
		nch.createChatContext(to_user_id, messageType, toUserData.WS, &map[string]interface{} {
			"type": "NEW_CHAT",
			"chat_id": chat_id,

		},a)
	}
}

func (nch *NewChatHandler) broadcastMessage(message types.MessageWS, messageType int, conn *websocket.Conn, actor types.Actor) {
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

func (nch *NewChatHandler) createChatContext(user_id int, messageType int, conn *websocket.Conn, response_data *map[string]interface{}, a Actor) {

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