package router

import (
	"main/controller"
	chat_router "main/router/chat"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine;
}

func NewRouter(router *gin.Engine) *Router {
	return &Router{Router: router}
}

func (r *Router) RegisterRouters(controller *controller.Controller) {
	api := r.Router.Group("/api");
	chatRouter := chat_router.NewChat(api);

	chatRouter.ChatRoutes(controller);
}