package controller

import (
	chat_controller "main/controller/chat"
	"main/database"
	websockets "main/websockets"
)

type Controller struct {
	Chat chat_controller.ChatController;
	WS websockets.WSConnection
}

func NewController(db *database.Database) *Controller{

	return &Controller{
		Chat: chat_controller.ChatController{DB: db},
		WS: websockets.WSConnection{DB: db},
	};
}