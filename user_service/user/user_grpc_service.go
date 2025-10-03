package user

import (
	"context"
	"log"
	"main/database/queries"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type Server struct {
    UnimplementedUserServiceServer
    users map[string]*UserResponse
    uq *queries.UserQuery
}

func NewServer(uq *queries.UserQuery) *Server {
    return &Server{
        users: make(map[string]*UserResponse),
        uq: uq,
    }
}

func (s *Server) GetUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
    var chatIdUserName map[string]string = make(map[string]string); 
    for chatId, userId  := range req.ChatIdAndUserIds {
        intId, _ := strconv.Atoi(userId);
        userData, err := s.uq.GetUserByID(intId);
        if err != nil {
            println(err)
            return nil, err
        }
        chatIdUserName[chatId] = userData.Name;
    }


    user := &UserResponse {ChatIdAndUserNames: chatIdUserName};
    
    return user, nil
}


func StartUserServer(uq *queries.UserQuery) {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    s := grpc.NewServer()
    RegisterUserServiceServer(s, NewServer(uq))
    
    log.Printf("UserService server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}