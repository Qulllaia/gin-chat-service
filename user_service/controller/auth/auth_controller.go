package auth_controller

import "github.com/gin-gonic/gin"

type AuthController struct {};

func (ac *AuthController) CreateUser(context *gin.Context) {
	context.JSON(200, gin.H{
		"done": true,
	})
}