package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"main/controller/dto"
	"main/database"
	"main/user"

	. "main/database/models"

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

	resultIds := make(map[string]string)
	var chats []dto.ChatListDTO
	for rows.Next() {
		var chat dto.ChatListDTO

		if err := rows.Scan(&chat.ID, &chat.Name, &chat.Chat_type_id, &chat.Users, &chat.Chat_background); err != nil {
			fmt.Println(err.Error())
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
			WHERE chat_id = $1)
		`, chat.ID)
		if err != nil {
			chat.LastMessage = "Сообщений пока не было!"
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
		`INSERT INTO "Chat" (users, name, chat_type_id, chat_background) 
		VALUES($1, $2, $3, NULL) RETURNING chat_id`,
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

func (cq *ChatQueries) DeleteChat(chatID int, userID int) ([]int64, error) {
	var users pq.Int64Array
	err := cq.DB.QueryRow(
		`SELECT users FROM "Chat" WHERE chat_id = $1`,
		chatID,
	).Scan(&users)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("chat not found")
		}
		return nil, err
	}

	isMember := false
	for _, id := range users {
		if int(id) == userID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, errors.New("user is not a member of this chat")
	}

	tx, err := cq.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`DELETE FROM "ChatHistory" WHERE chat_id = $1`, chatID); err != nil {
		return nil, err
	}

	if _, err = tx.Exec(`DELETE FROM "Chat" WHERE chat_id = $1`, chatID); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return []int64(users), nil
}

func (cq *ChatQueries) GetChatMembers(chatID int, requesterID int) ([]dto.ChatMemberDTO, error) {
	var users pq.Int64Array
	var chatTypeID string
	err := cq.DB.QueryRow(
		`SELECT users, chat_type_id FROM "Chat" WHERE chat_id = $1`,
		chatID,
	).Scan(&users, &chatTypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("chat not found")
		}
		return nil, err
	}

	if chatTypeID != "GROUPCHAT" {
		return nil, errors.New("members list is available only for group chats")
	}

	isMember := false
	for _, id := range users {
		if int(id) == requesterID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, errors.New("user is not a member of this chat")
	}

	intIDs := make([]int, 0, len(users))
	for _, id := range users {
		intIDs = append(intIDs, int(id))
	}

	names, err := cq.grpc.GetUserNamesByIDs(intIDs)
	if err != nil {
		return nil, err
	}

	members := make([]dto.ChatMemberDTO, 0, len(users))
	for _, id := range users {
		name := names[int(id)]
		if name == "" {
			name = "Пользователь"
		}
		members = append(members, dto.ChatMemberDTO{
			ID:   id,
			Name: name,
		})
	}

	return members, nil
}
