package controller

import (
	chat_controller "main/controller/chat"
)

type Controller struct {
	Chat chat_controller.ChatController;
}

func NewController() *Controller{

	return &Controller{
		Chat: chat_controller.ChatController{},
	};
}