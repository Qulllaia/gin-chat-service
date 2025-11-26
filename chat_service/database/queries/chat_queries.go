package queries

import (
	"strconv"

	"main/controller/dto"
	"main/database"
	. "main/database/models"
	"main/user"

	"github.com/lib/pq"
)

type ChatQueries struct {
	grpc *user.Server
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
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.Id, &message.Message, &message.Chat_id, &message.User_id, &message.Timestamp); err != nil {
			return nil, err
		}
		message.IsThisUserSender = current_user_id == message.User_id
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (cq *ChatQueries) GetUsersChats(currentId int, users *[]dto.ChatListDTO) error {
	rows, err := cq.DB.Query(`
        SELECT chat_id, name, chat_type_id, users, chat_background
        FROM "Chat" c
        WHERE $1 = ANY(users)
    `, currentId)
	if err != nil {
		println(err.Error())
		return err
	}
	defer rows.Close()

	var resultIds map[string]string = make(map[string]string)
	var chats []dto.ChatListDTO
	for rows.Next() {
		var chat dto.ChatListDTO

		if err := rows.Scan(&chat.ID, &chat.Name, &chat.Chat_type_id, &chat.Users, &chat.Chat_background); err != nil {
			return err
		}

		if chat.Chat_type_id == "PRIVATECHAT" {
			for _, i := range chat.Users {
				if i != int64(currentId) {
					resultIds[strconv.Itoa(chat.ID)] = strconv.Itoa(int(i))
					chat.User_id = &i
				}
			}
		}

		err = cq.DB.Get(&chat.LastMessage,
			`SELECT ch.message FROM "ChatHistory" ch 
			WHERE ch.id in (select max(id) as message_max_id FROM "ChatHistory" 
			WHERE chat_id = $1)`,
			chat.ID)
		if err != nil {
			panic(err)
		}

		chats = append(chats, chat)
	}

	if len(resultIds) > 0 {
		cq.grpc.GetUserInfo(resultIds, &chats)
		*users = chats
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func (cq *ChatQueries) CreateMultipleUserChat(ids []int64, groupName string) (error, int64) {
	var resultChatId int64
	userIDsArray := pq.Array(ids)
	err := cq.DB.QueryRow(
		`INSERT INTO "Chat" (users, name, chat_type_id) 
		VALUES($1, $2, $3) RETURNING chat_id`,
		userIDsArray, groupName, "GROUPCHAT",
	).Scan(&resultChatId)

	return err, resultChatId
}

func (cq *ChatQueries) AddBachgroundToChat(chat_id int, chat_background string) error {
	_, err := cq.DB.Exec(
		`UPDATE "Chat" set chat_background = $1 where chat_id = $2`,
		chat_background, chat_id,
	)
	if err != nil {
		return err
	}

	return nil
}
