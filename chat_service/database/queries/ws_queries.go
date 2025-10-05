package queries

import (
	"main/database"
	. "main/database/models"
	"time"

	"github.com/lib/pq"
)

type WSQueries struct {
	*database.Database
}

func WSQueryConstructor(db *database.Database) *WSQueries {
	return &WSQueries{db}
}


func (wsq *WSQueries) GetUserChatList(user_id int, user *User) (error) {
    rows, err := wsq.DB.Query(`
        SELECT chat_id 
        FROM "Chat"
        WHERE $1 = ANY(users) 
    `, user_id);

    var chatIDs []int64

    for rows.Next() {
        var chatID int64 
        err := rows.Scan(&chatID) 
        if err != nil {
            return err
        }
        chatIDs = append(chatIDs, chatID)
    }

    if err := rows.Err(); err != nil {
        return err
    }

    user.Chat_list = chatIDs

	return err;
}

func (wsq *WSQueries) InsertMessageIntoChatHistory(chat_id, user_id int, message string) (error) {
	_, err := wsq.DB.Exec(
		`INSERT INTO "ChatHistory" (message, chat_id, user_id, timestamp) 
		VALUES($1, $2, $3, $4)`, 
		message, chat_id, user_id, time.Now(),
		);
	return err;
}

func (wsq *WSQueries) CreateChatWithMessage(user_ids ...int) (error, int) {
	var chat_id int;
	userIDsArray := pq.Array(user_ids)
	err := wsq.DB.QueryRow(
		`INSERT INTO "Chat" (users, name, chat_type) 
		VALUES($1, $2, $3) RETURNING chat_id`, 
		userIDsArray, nil, "PRIVATECHAT",
	).Scan(&chat_id);

	return err, chat_id;
}