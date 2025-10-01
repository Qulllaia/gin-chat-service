package user

import (
	"context"
	"strconv"
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


    // res, err := userClient.GetUser(ctx,  &UserRequest{UserId: "34"})
    // if err != nil {
    //     println(res)
    //     return nil, err
    // }

    // println(res.UserId);
    // println(res.state);

    return &Server{
        userClient: userClient,
    }, nil
}

func (s *Server) GetUserInfo(id int) (*UserResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    
    defer cancel()

    res, err  := s.userClient.GetUser(ctx, &UserRequest{UserId: strconv.Itoa(id)})

    if err != nil {
        println(res)
        return nil, err
    }

    // println(res.UserId);

    return res, nil;
}