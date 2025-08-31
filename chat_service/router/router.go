package router

import (
	"main/config"
	"main/controller"
	"main/middleware"
	chat_router "main/router/chat"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine;
}

func NewRouter(router *gin.Engine) *Router {
	return &Router{Router: router}
}

func (r *Router) RegisterRouters(controller *controller.Controller, config *config.Config) {
	api := r.Router.Group("/api", middleware.ConfigMiddleware(config));
	chatRouter := chat_router.NewChat(api);

	chatRouter.ChatRoutes(controller);
}