package queries

import (
	"main/database"
	. "main/database/models"
)

type ChatQueries struct {
	*database.Database
}

func ChatQueryConstructor(db *database.Database) *ChatQueries {
	return &ChatQueries{db}
}

func (cq *ChatQueries) GetMessageHistory(current_user_id int64) ([]Message, error) {
    rows, err := cq.DB.Query(`
        SELECT id, message, chat_id, user_id, "timestamp"
        FROM "ChatHistory"  ORDER BY "timestamp" DESC
    `)
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