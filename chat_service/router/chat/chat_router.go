package chat_router

import (
	"main/controller"
	"main/controller/dto"
	"main/database/models"
	"main/middleware"
	"main/types"

	"github.com/gin-gonic/gin"
)

type Chat struct {
	ChatRouter *gin.RouterGroup
}

func NewChat(router *gin.RouterGroup) *Chat {
	return &Chat{ChatRouter: router}
}

func (a *Chat) ChatRoutes(controller *controller.Controller) {
	api := a.ChatRouter.Group("chat", middleware.AuthMiddleware)
	{
		api.GET("/history/:id", middleware.ErrorMiddleware[[]models.Message](controller.Chat.GetHistoryList))
		api.GET("/ws", controller.WS.WebsocketsInit)
		api.GET("/chats", middleware.ErrorMiddleware[[]dto.ChatListDTO](controller.Chat.GetUsersChats))
		api.POST("/chats", middleware.ErrorMiddleware[*int64](controller.Chat.CreateChatWithMultipleUsers))
		api.POST("/background", middleware.ErrorMiddleware[*types.ImageResponse](controller.Chat.SetBackGround))
	}
}

