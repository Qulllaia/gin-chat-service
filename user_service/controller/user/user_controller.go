package user_controller

import "github.com/gin-gonic/gin"

type UserController struct{};

func (uc *UserController) CreateUser(context *gin.Context) {
	context.JSON(200, gin.H{
		"done": true,
	})
}