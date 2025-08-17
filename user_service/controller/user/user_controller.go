package user_controller

import (
	"main/database/queries"

	"github.com/gin-gonic/gin"
)

type UserController struct{
	UQ *queries.UserQuery
};

func (uc *UserController) CreateUser(context *gin.Context) {

	err := uc.UQ.InsertUser("lol", "kek");
	if err != nil {
		context.JSON(400, gin.H{
			"error": "createUserException",
			"message": err.Error(),
		})
	} else {
		context.JSON(200, gin.H{
			"done": true,
		})
	}
}