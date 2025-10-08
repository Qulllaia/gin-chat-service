package chat_controller

import (
	"main/controller/utils"
	"main/database/queries"
	"net/http"

	. "main/controller/dto"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	CQ *queries.ChatQueries
}

func (cc *ChatController) GetHistoryList(context *gin.Context) {

	claims, err := utils.ExtractClaimsFromCookie(context);

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"done": false,
			"message": err.Error(),
		})
		return;
	}
	var chatID ChatIDURI;

	if err := context.ShouldBindUri(&chatID); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messages, err := cc.CQ.GetMessageHistory(int64(claims.UserID), int64(chatID.ID));

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


func (cc *ChatController) GetUsersChats(context *gin.Context) {

	claims, err := utils.ExtractClaimsFromCookie(context);

	users, err := cc.CQ.GetUsersChats(int(claims.UserID));
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "GetUsersChatsException",
			"message": err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"done": true,
			"result": users,
		})
	}
}

func (cc *ChatController) CreateChatWithMultipleUsers(context *gin.Context) {
	claims, err := utils.ExtractClaimsFromCookie(context);
	
	var idsJson UsersIDList;
	
	if err = context.ShouldBindBodyWithJSON(&idsJson); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "CreateChatWithMultipleUsers",
			"message": err.Error(),
		})
	}

	err, resultId := cc.CQ.CreateMultipleUserChat(append(idsJson.IDs, int64(claims.UserID)), idsJson.GroupName)

	context.JSON(http.StatusCreated, gin.H{
		"done": true,
		"result": resultId,
	})
	
}