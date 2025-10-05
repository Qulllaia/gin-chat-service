package queries

import (
	"main/controller/dto"
	"main/database"
	. "main/database/models"
	"main/user"
	"strconv"
)

type ChatQueries struct {
    *user.Server
	*database.Database
}

func ChatQueryConstructor(db *database.Database, gprc *user.Server) *ChatQueries {
	return &ChatQueries{gprc, db}
}

func (cq *ChatQueries) GetMessageHistory(current_user_id, chat_id int64) ([]Message, error) {
    rows, err := cq.DB.Query(`
        SELECT id, message, chat_id, user_id, "timestamp"
        FROM "ChatHistory" WHERE chat_id = $1 ORDER BY "timestamp" ASC
    `, chat_id)
    if err != nil {
		return nil, err;
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var message Message
        if err := rows.Scan(&message.Id, &message.Message, &message.Chat_id, &message.User_id, &message.Timestamp); err != nil {
			return nil, err;
        }
		message.IsThisUserSender = current_user_id == message.User_id;
        messages = append(messages, message)
    }
    
    if err = rows.Err(); err != nil {
		return nil, err;
	}

	return messages, nil;
}

func (cq *ChatQueries) GetUsersChats(currentId int) ([]dto.ChatListDTO, error) {
    rows, err := cq.DB.Query(`
        SELECT chat_id, chat_type, users
        FROM "Chat" c
        WHERE $1 = ANY(users)
    `, currentId)
    if err != nil {
        println(err.Error());
        return nil, err
    }
    defer rows.Close()
    
    var resultIds map[string]string = make(map[string]string);
    var chats []dto.ChatListDTO
    for rows.Next() {
        var chat dto.ChatListDTO
        if err := rows.Scan(&chat.ID, &chat.Chat_type, &chat.Users); err != nil {
            return nil, err
        }
        if chat.Chat_type == "PRIVATECHAT" {
            for _, i := range chat.Users {
                if i != int64(currentId) {
                    resultIds[strconv.Itoa(chat.ID)] = strconv.Itoa(int(i))    
                } 
            }
            
        }
        chats = append(chats, chat)
    }

    if len(resultIds) > 0 {

        userGRPCResponse, err := cq.GetUserInfo(resultIds);
        if err != nil {
            println("Fatal Error to Get User Info with GRPC")
            println(err.Error())
            return nil, err
        } else {
            
            for index, i := range chats {
                chatId := strconv.Itoa(i.ID);
                val, exists := userGRPCResponse.ChatIdAndUserNames[chatId];
                
                if exists {
                    chats[index].Name = &val; 
                }

            }
        }
        
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return chats, nil
}