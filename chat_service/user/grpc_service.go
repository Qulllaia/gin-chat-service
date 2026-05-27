package user

import (
	"context"
	"strconv"
	"time"

	"main/controller/dto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	userClient UserServiceClient
}

func ConnectSerivce(userServiceAddr string) (*Server, error) {
	conn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	userClient := NewUserServiceClient(conn)

	return &Server{
		userClient: userClient,
	}, nil
}

func (s *Server) GetUserInfo(chatIdAndUserIds map[string]string, chats *[]dto.ChatListDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	res, err := s.userClient.GetUser(ctx, &UserRequest{ChatIdAndUserIds: chatIdAndUserIds})
	if err != nil {
		println(res)
		return err
	}

	for index, i := range *chats {
		chatId := strconv.Itoa(i.ID)
		val, exists := res.ChatIdAndUserNames[chatId]

		if exists {
			user_id, _ := strconv.Atoi(chatIdAndUserIds[chatId])
			user_id64 := int64(user_id)
			(*chats)[index].User_id = &user_id64
			(*chats)[index].Name = &val
		}

	}
	return nil
}

func (s *Server) GetUserNamesByIDs(userIDs []int) (map[int]string, error) {
	if len(userIDs) == 0 {
		return map[int]string{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqMap := make(map[string]string, len(userIDs))
	for _, id := range userIDs {
		key := strconv.Itoa(id)
		reqMap[key] = key
	}

	res, err := s.userClient.GetUser(ctx, &UserRequest{ChatIdAndUserIds: reqMap})
	if err != nil {
		return nil, err
	}

	names := make(map[int]string, len(res.ChatIdAndUserNames))
	for chatID, name := range res.ChatIdAndUserNames {
		id, err := strconv.Atoi(chatID)
		if err != nil {
			continue
		}
		names[id] = name
	}

	return names, nil
}

