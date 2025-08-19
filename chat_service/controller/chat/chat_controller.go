package chat_controller

import (
	"main/database"
	. "main/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	DB *database.Database
}

func (cc *ChatController) GetHistoryList(context *gin.Context) {
    rows, err := cc.DB.DB.Query(`
        SELECT id, message, chat_id, user_id, "timestamp"
        FROM "ChatHistory"  ORDER BY "timestamp" DESC
    `)
    if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"done": false,
			"result": err.Error(),
		})
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var message Message
        if err := rows.Scan(&message.Id, &message.Message, &message.Chat_id, &message.User_id, &message.Timestamp); err != nil {
			context.JSON(http.StatusOK, gin.H{
				"done": false,
				"result": err.Error(),
			})
        }
        messages = append(messages, message)
    }
    
    if err = rows.Err(); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"done": false,
			"result": err.Error(),
		})
	}

	context.JSON(http.StatusOK, gin.H{
		"done": true,
		"result": messages,
	})
}