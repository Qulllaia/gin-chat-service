package queries

import (
	"main/controller/dto"
	"main/database"
	. "main/database/models"
	"main/user"
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

        // SELECT chat_id,

        //     CASE 
        //         WHEN c.chat_type = 'PRIVATECHAT' THEN u.name 
        //         ELSE c.name 
        //     END AS name

        // FROM "Chat" c
        // LEFT JOIN "user" u on u.id = (SELECT unnest(c.users) EXCEPT SELECT $1)
        // WHERE $1 = ANY(users) 

func (cq *ChatQueries) GetUsersChats(id int) ([]dto.ChatListDTO, error) {
    rows, err := cq.DB.Query(`
        SELECT chat_id, chat_type, users
        FROM "Chat" c
        WHERE $1 = ANY(users)
    `, id)
    if err != nil {
        println(err.Error());
        return nil, err
    }
    defer rows.Close()

    var chats []dto.ChatListDTO
    for rows.Next() {
        var chat dto.ChatListDTO
        if err := rows.Scan(&chat.ID, &chat.Chat_type, &chat.Users); err != nil {
            return nil, err
        }
        // println(chat.ID);
        if chat.Chat_type == "PRIVATECHAT" {
            for _, i := range chat.Users {
                if i != int64(id) {

                    userGRPCResponse, err := cq.GetUserInfo(int(i));
                    if err != nil {
                        println("Fatal Error to Get User Info with GRPC")
                        println(err.Error())
                        return nil, err
                    } else {
                        // println(userGRPCResponse.UserId)
                        chat.Name = &userGRPCResponse.Name;
                    }
                } 
            }
        }
        chats = append(chats, chat)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return chats, nil
}