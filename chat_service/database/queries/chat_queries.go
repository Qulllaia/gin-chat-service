package queries

import (
	"main/controller/dto"
	"main/database"
	. "main/database/models"
)

type ChatQueries struct {
	*database.Database
}

func ChatQueryConstructor(db *database.Database) *ChatQueries {
	return &ChatQueries{db}
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

func (uq *ChatQueries) GetUsersChats(id int) ([]dto.ChatListDTO, error) {
    rows, err := uq.DB.Query(`
        SELECT chat_id, name
        FROM "Chat" 
        WHERE $1 = ANY(users) 
    `, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var chats []dto.ChatListDTO
    for rows.Next() {
        var chat dto.ChatListDTO
        if err := rows.Scan(&chat.ID, &chat.Name); err != nil {
            return nil, err
        }
        chats = append(chats, chat)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return chats, nil
}