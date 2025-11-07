package auth_router

import (
	"main/controller"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	AuthRouter *gin.RouterGroup
}

func NewAuth(router *gin.RouterGroup) *Auth {
	return &Auth{AuthRouter: router}
}

func (a *Auth) AuthRoutes(controller *controller.Controller) {
	api := a.AuthRouter.Group("auth")
	{
		// api.POST("/reg", controller.Auth.RegisterUser)
		api.POST("/login", controller.Auth.LoginUser)
		api.POST("/verify", controller.Auth.SMTPApprove)
		api.GET("/verify/:token", controller.Auth.VerifyResult)
	}
}
