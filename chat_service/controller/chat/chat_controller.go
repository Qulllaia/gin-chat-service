package chat_controller

import (
	"encoding/json"
	"fmt"
	"main/controller/dto"
	"main/controller/utils"
	"main/database/queries"
	"main/redis"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "main/controller/dto"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	CQ *queries.ChatQueries
	RDB *redis.RedisConnector
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
	
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "GetUsersChatsException",
			"message": err.Error(),
		})
	}

	var users []dto.ChatListDTO; 
	if result, _ := cc.RDB.DoesDataExists(strconv.Itoa(claims.UserID)); *result != 1  {	
		err := cc.CQ.GetUsersChats(int(claims.UserID), &users);
		jsonData, err := json.Marshal(users);
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "GetUsersChatsException",
				"message": err.Error(),
			})
		}
		fmt.Println(users)
		err = cc.RDB.SetData(strconv.Itoa(claims.UserID), string(jsonData))
	} else {
		stringUsers, err := cc.RDB.GetData(strconv.Itoa(claims.UserID))

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "GetUsersChatsException",
				"message": err.Error(),
			})
		}


		err = json.Unmarshal([]byte(stringUsers), &users)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "GetUsersChatsException",
				"message": err.Error(),
			})
		}

	}

	fmt.Println(users)
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

	fullUserList := append(idsJson.IDs, int64(claims.UserID))
	err, resultId := cc.CQ.CreateMultipleUserChat(fullUserList, idsJson.GroupName)

	for _, i := range(fullUserList) {	
		err = cc.RDB.DeleteData(strconv.Itoa(int(i)));
	}


	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "CreateChatWithMultipleUsers",
			"message": err.Error(),
		})
	}	

	context.JSON(http.StatusCreated, gin.H{
		"done": true,
		"result": resultId,
	})		
}

func (cc *ChatController) SetBackGround(context *gin.Context) {
	file, err := context.FormFile("image");
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
	}
	
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	
	ext := filepath.Ext(file.Filename)
	if !allowedTypes[ext] {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type",
		})
		return
	}
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s", timestamp, filepath.Base(file.Filename))

	filename = strings.ReplaceAll(filename, " ", "_")
    filename = strings.ReplaceAll(filename, ":", "-")
    
	filePath := filepath.Join("./static/background", filename)
	
	if err := context.SaveUploadedFile(file, filePath); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	imageURL := fmt.Sprintf("/background/%s", filename)
	
	context.JSON(http.StatusOK, gin.H{
		"message":   "Image uploaded successfully",
		"filename":  filename,
		"url":       imageURL,
		"full_url":  context.Request.Host + imageURL,
	})
}