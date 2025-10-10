package chat_router

import (
	"main/controller"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

type Chat struct {
	ChatRouter *gin.RouterGroup
}

func NewChat(router *gin.RouterGroup) *Chat {
	return &Chat{ ChatRouter: router};
}

func (a *Chat) ChatRoutes(controller *controller.Controller) {
	api := a.ChatRouter.Group("chat", middleware.AuthMiddleware)
	{
		api.GET("/history/:id", controller.Chat.GetHistoryList)
		api.GET("/ws", controller.WS.WebsocketsInit)
		api.GET("/chats", controller.Chat.GetUsersChats)
		api.POST("chats", controller.Chat.CreateChatWithMultipleUsers)
		api.POST("background", controller.Chat.SetBackGround)
	}
}