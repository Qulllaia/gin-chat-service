package chat_router

import (
	"main/controller"

	"github.com/gin-gonic/gin"
)

type Chat struct {
	ChatRouter *gin.RouterGroup
}

func NewChat(router *gin.RouterGroup) *Chat {
	return &Chat{ ChatRouter: router};
}

func (a *Chat) ChatRoutes(controller *controller.Controller) {
	api := a.ChatRouter.Group("chat")
	{
		api.GET("/history", controller.Chat.GetHistoryList)
		api.GET("/ws", controller.WS.WebsocketsInit)
	}
}