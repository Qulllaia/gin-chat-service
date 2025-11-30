package chat_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"main/controller/dto"
	. "main/controller/dto"
	"main/database/models"
	"main/database/queries"
	"main/redis"
	. "main/types"
	"main/utils"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	CQ  *queries.ChatQueries
	RDB *redis.RedisConnector
}

func (cc *ChatController) GetHistoryList(context *gin.Context) (ErrorType, []models.Message, error) {
	claims, err := utils.ExtractClaimsFromCookie(context)
	if err != nil {
		return ClaimsExtractingError, nil, err
	}

	var chatID ChatIDURI
	if err := context.ShouldBindUri(&chatID); err != nil {
		return UriParsingError, nil, err
	}

	messages, err := cc.CQ.GetMessageHistory(int64(claims.UserID), int64(chatID.ID))
	if err != nil {
		return DatabaseError, nil, err
	}

	return NoError, messages, nil
}

func (cc *ChatController) GetUsersChats(context *gin.Context) (ErrorType, []dto.ChatListDTO, error) {
	claims, err := utils.ExtractClaimsFromCookie(context)
	if err != nil {
		return ClaimsExtractingError, nil, err
	}
	stringUserId := strconv.Itoa(claims.UserID)
	var users []dto.ChatListDTO
	if result, _ := cc.RDB.DoesDataExists(stringUserId); *result != 1 {

		if err := cc.CQ.GetUsersChats(int(claims.UserID), &users); err != nil {
			return CacheError, nil, err
		}

		jsonData, err := json.Marshal(users)
		if err != nil {
			return JsonParsingError, nil, err
		}

		if err := cc.RDB.SetData(stringUserId, string(jsonData)); err != nil {
			return CacheError, nil, err
		}

	} else {
		stringUsers, err := cc.RDB.GetData(stringUserId)
		if err != nil {
			return CacheError, nil, err
		}

		if err = json.Unmarshal([]byte(stringUsers), &users); err != nil {
			return JsonParsingError, nil, err
		}

	}

	return NoError, users, nil
}

func (cc *ChatController) CreateChatWithMultipleUsers(context *gin.Context) (ErrorType, *int64, error) {
	claims, err := utils.ExtractClaimsFromCookie(context)
	if err != nil {
		return ClaimsExtractingError, nil, err
	}

	var idsJson UsersIDList
	if err := context.ShouldBindBodyWithJSON(&idsJson); err != nil {
		return JsonParsingError, nil, err
	}

	fullUserList := append(idsJson.IDs, int64(claims.UserID))

	err, resultId := cc.CQ.CreateMultipleUserChat(fullUserList, idsJson.GroupName)
	if err != nil {
		return DatabaseError, nil, err
	}

	for _, i := range fullUserList {
		if err := cc.RDB.DeleteData(strconv.Itoa(int(i))); err != nil {
			return CacheError, nil, err
		}
	}

	return NoError, &resultId, nil
}

func (cc *ChatController) SetBackGround(context *gin.Context) (ErrorType, *ImageResponse, error) {
	claims, err := utils.ExtractClaimsFromCookie(context)
	if err != nil {
		return ClaimsExtractingError, nil, err
	}

	file, err := context.FormFile("image")
	chat_id := context.PostForm("chat_id")

	if err != nil {
		return FileExtractingError, nil, err
	}

	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := filepath.Ext(file.Filename)
	if !allowedTypes[ext] {
		return FileTypeError, nil, errors.New("Invalid file type")
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s", timestamp, filepath.Base(file.Filename))

	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, ":", "-")

	filePath := filepath.Join("./static/background", filename)

	if err := context.SaveUploadedFile(file, filePath); err != nil {
		return FileSavingError, nil, err
	}

	imageURL := fmt.Sprintf("/background/%s", filename)

	intChat_id, err := strconv.Atoi(chat_id)
	if err != nil {
		return ConversionError, nil, err
	}

	if err = cc.CQ.AddBachgroundToChat(intChat_id, imageURL); err != nil {
		return DatabaseError, nil, err
	}

	if err = cc.RDB.DeleteData(strconv.Itoa(claims.UserID)); err != nil {
		return CacheError, nil, err
	}

	return NoError, &ImageResponse{
		Message:  "Image uploaded successfully",
		Filename: filename,
		Url:      imageURL,
		Full_url: context.Request.Host + imageURL,
	}, nil
}

