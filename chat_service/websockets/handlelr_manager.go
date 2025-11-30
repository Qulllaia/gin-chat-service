package websockets

import (
	. "main/types"
	"main/websockets/handlers"
)

type HandlerManager struct {
	handlers map[MessageType]Handler
}

func NewHandlerManager(actor *ConnectorActor) *HandlerManager {
	messageHandler := handlers.NewMessageHandler()
	newChatHandler := handlers.NewNewChatHandler()
	newMultiChatHandler := handlers.NewNewMultiChatHandler()
	newOnlineStatusHandler := handlers.NewOnlineStatusHandler()
	newMediaHandler := handlers.NewMediaHandler()
	handlersMap := make(map[MessageType]Handler)

	handlersMap[NEW_CHAT] = newChatHandler
	handlersMap[MESSAGE] = messageHandler
	handlersMap[NEW_MULTIPLE_CHAT] = newMultiChatHandler
	handlersMap[USER_STATUS] = newOnlineStatusHandler
	handlersMap[MEDIA] = newMediaHandler
	return &HandlerManager{handlers: handlersMap}
}

