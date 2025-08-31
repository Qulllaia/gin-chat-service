package controller

import (
	chat_controller "main/controller/chat"
	"main/database"
	"main/database/queries"
	websockets "main/websockets"
)

type Controller struct {
	Chat chat_controller.ChatController;
	WS websockets.WSConnection
}

func NewController(db *database.Database, cq *queries.ChatQueries, a *websockets.ConnectorActor) *Controller{

	return &Controller{
		Chat: chat_controller.ChatController{CQ: cq},
		WS: websockets.WSConnection{Actor: a},
	};
}