package user

import (
	"context"
	"time"

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

func (s *Server) GetUserInfo(chatIdAndUserIds map[string]string) (*UserResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    
    defer cancel()

    res, err  := s.userClient.GetUser(ctx, &UserRequest{ChatIdAndUserIds: chatIdAndUserIds})

    if err != nil {
        println(res)
        return nil, err
    }

    return res, nil;
}