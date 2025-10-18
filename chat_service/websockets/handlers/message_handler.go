package handlers

import (
	"main/types"
	"main/websockets/handlers/base"

	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	*base.BaseHandler
}

func NewMessageHandler() types.Handler {
	return &MessageHandler{}
}

func (mh *MessageHandler) Handle(message types.MessageWS, messageType int, conn *websocket.Conn, actor types.Actor) {
	mh.BroadcastMessage(message, messageType, conn, actor)
}