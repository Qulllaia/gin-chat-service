package websockets

import (
	. "main/types"
	"main/websockets/handlers"
)

type HandlerManager struct {
	handlers map[string]Handler
}

func NewHandlerManager(actor *ConnectorActor) *HandlerManager {
	messageHandler := handlers.NewMessageHandler()
	newChatHandler := handlers.NewNewChatHandler()
	newMultiChatHandler := handlers.NewNewMultiChatHandler()	
	newOnlineStatusHandler := handlers.NewOnlineStatusHandler()	

	handlersMap := make(map[string]Handler)

	handlersMap["NEW_CHAT"] = newChatHandler
	handlersMap["MESSAGE"] = messageHandler
	handlersMap["NEW_MULTIPLE_CHAT"] = newMultiChatHandler
	handlersMap["USER_STATUS"] = newOnlineStatusHandler
	return &HandlerManager{ handlers: handlersMap}
}