package queries

import (
	"main/controller/utils"
	"main/database"
	. "main/database/models"
	"time"
)

type WSQueries struct {
	*database.Database
}

func WSQueryConstructor(db *database.Database) *WSQueries {
	return &WSQueries{db}
}


func (wsq *WSQueries) GetUserChatList(claims *utils.Claims, user *User) (error) {
    err := wsq.DB.QueryRow(`
        SELECT id, chat_list 
        FROM "user" 
        WHERE id = $1
    `, claims.UserID).Scan(&user.ID, &user.Chat_list)

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