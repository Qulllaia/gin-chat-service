package chat_controller

import (
	"main/controller/utils"
	"main/database/queries"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	CQ *queries.ChatQueries
}

func (cc *ChatController) GetHistoryList(context *gin.Context) {

	cookie := context.Request.Cookies();

	jwt_token := "";

	for _, val := range cookie {
		if val.Name == "session_token" {
			jwt_token = val.Value;
		}
	}

	claims, err := utils.DecodeJWT(jwt_token);

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"done": false,
			"message": err.Error(),
		})
		return;
	}

	messages, err := cc.CQ.GetMessageHistory(int64(claims.UserID));

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"done": false,
			"message": err.Error(),
		})
		return;
	}

	context.JSON(http.StatusOK, gin.H{
		"done": true,
		"result": messages,
	})
}