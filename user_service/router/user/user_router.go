package user_router

import (
	"main/controller"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserRouter *gin.RouterGroup
}

func NewUser(router *gin.RouterGroup) *User {
	return &User{ UserRouter: router};
}

func (a *User) UserRoutes(controller *controller.Controller) {
	api := a.UserRouter.Group("user")
	{
		api.GET("/create", controller.User.CreateUser)
		api.GET("/get")
		api.GET("/delete")
		api.GET("/update")
	}
}