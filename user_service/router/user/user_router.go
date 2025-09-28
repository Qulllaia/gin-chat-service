package user_router

import (
	"main/controller"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserRouter *gin.RouterGroup
}



func NewUser(router *gin.RouterGroup) *User {
	return &User{ UserRouter: router};
}






func (a *User) UserRoutes(controller *controller.Controller) {
	api := a.UserRouter.Group("user", middleware.AuthMiddleware)
	{
		api.POST("/create", controller.User.CreateUser)
		api.GET("/get/:id", controller.User.GetUserByID)
		api.GET("/get/except", controller.User.GetUserExceptCurrent)
		api.GET("/get", controller.User.GetAllUsers)
		api.DELETE("/delete/:id", controller.User.DeleteUser)
		api.PUT("/update", controller.User.UpdateUser)
	}
}