package auth_router

import (
	"main/controller"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	AuthRouter *gin.RouterGroup
}

func NewAuth(router *gin.RouterGroup) *Auth {
	return &Auth{ AuthRouter: router};
}

func (a *Auth) AuthRoutes(controller *controller.Controller) {
	api := a.AuthRouter.Group("auth")
	{
		api.GET("/reg", controller.Auth.CreateUser)
		api.GET("/login")
	}
}