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
    // Подключение к UserService
    conn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }
    
    userClient := NewUserServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    res, err := userClient.GetUser(ctx,  &UserRequest{UserId: "34"})
    if err != nil {
        println(res)
    }

    println(res.UserId);
    println(res.state);

    return &Server{
        userClient: userClient,
    }, nil
}
