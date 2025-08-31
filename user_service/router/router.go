package router

import (
	"main/config"
	"main/controller"
	"main/middleware"
	auth_router "main/router/auth"
	user_router "main/router/user"

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
	authRouter := auth_router.NewAuth(api);
	userRouter := user_router.NewUser(api);

	authRouter.AuthRoutes(controller);
	userRouter.UserRoutes(controller);
}