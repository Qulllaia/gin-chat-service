package handlers

import (
	"encoding/json"
	"fmt"
	"main/types"
	"main/utils"
	"main/websockets/handlers/base"

	"github.com/gorilla/websocket"
)

type OnlineStatusHandler struct {
	*base.BaseHandler
}

func NewOnlineStatusHandler() types.Handler {
	return &OnlineStatusHandler{}
}

func (osh *OnlineStatusHandler) Handle(message types.MessageWS, messageType int, conn *websocket.Conn, actor types.Actor) {
	fmt.Println(message)	
	osh.ThoseUsersThatNeedToUpdateStatus(message, messageType, conn, actor)
}

// TODO: Привести входные данные к общему формату
func (osh *OnlineStatusHandler) ThoseUsersThatNeedToUpdateStatus(message types.MessageWS, messageType int, conn *websocket.Conn, actor types.Actor) {
	connsToUsers := actor.GetConnectoinsToUsers()
	usersToWSData := actor.GetUserConnections()

	userIdThatChangeStatus := connsToUsers[conn]
	userThatChangeStatus := usersToWSData[userIdThatChangeStatus];

	chats := userThatChangeStatus.Chat_ids;

	for id, data := range usersToWSData {
		for _, chatId := range chats {
			if utils.Contains(data.Chat_ids, int64(chatId)) {
				jsonData, err := json.Marshal(&map[string]interface{}{
					message.Message: userIdThatChangeStatus,
				})
				if err != nil {
					fmt.Println(err)
				}
				if err := usersToWSData[id].WS.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					fmt.Println(err)
				}
				break
			} 
		}	
	}
}