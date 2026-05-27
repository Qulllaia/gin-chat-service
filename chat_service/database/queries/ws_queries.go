package queries

import (
	"database/sql"
	"errors"
	"time"

	"main/database"
	. "main/database/models"

	"github.com/lib/pq"
)

type WSQueries struct {
	*database.Database
}

func WSQueryConstructor(db *database.Database) *WSQueries {
	return &WSQueries{db}
}

func (wsq *WSQueries) GetUserChatList(user_id int, user *User) error {
	rows, err := wsq.DB.Query(`
        SELECT chat_id 
        FROM "Chat"
        WHERE $1 = ANY(users) 
    `, user_id)

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

	return err
}

func (wsq *WSQueries) InsertMessageIntoChatHistory(chat_id, user_id int, message string) error {
	_, err := wsq.DB.Exec(
		`INSERT INTO "ChatHistory" (message, chat_id, user_id, timestamp) 
		VALUES($1, $2, $3, $4)`,
		message, chat_id, user_id, time.Now(),
	)
	return err
}

func (wsq *WSQueries) FindPrivateChatBetweenUsers(userID1, userID2 int) (int, bool, error) {
	var chatID int
	err := wsq.DB.QueryRow(`
		SELECT chat_id
		FROM "Chat"
		WHERE chat_type_id = 'PRIVATECHAT'
		  AND $1 = ANY(users)
		  AND $2 = ANY(users)
		  AND cardinality(users) = 2
		LIMIT 1
	`, userID1, userID2).Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}
	return chatID, true, nil
}

func (wsq *WSQueries) CreateChatWithMessage(user_ids ...int) (error, int) {
	if len(user_ids) != 2 {
		return errors.New("private chat requires exactly two users"), 0
	}

	if chatID, found, err := wsq.FindPrivateChatBetweenUsers(user_ids[0], user_ids[1]); err != nil {
		return err, 0
	} else if found {
		return nil, chatID
	}

	return wsq.createPrivateChat(user_ids)
}

func (wsq *WSQueries) createPrivateChat(user_ids []int) (error, int) {
	var chat_id int
	userIDsArray := pq.Array(user_ids)
	err := wsq.DB.QueryRow(
		`INSERT INTO "Chat" (users, name, chat_type_id, chat_background) 
		VALUES($1, $2, $3, NULL) RETURNING chat_id`,
		userIDsArray, nil, "PRIVATECHAT",
	).Scan(&chat_id)

	return err, chat_id
}

