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
    log.Printf("Received GetUser request for user_id: %s", req.UserId)
    
    // user, exists := s.users[req.UserId]

    intId, _ := strconv.Atoi(req.UserId);

    userData, err := s.uq.GetUserByID(intId);

    if err != nil {
        println(err)
        return nil, err
    }

    user := &UserResponse {UserId: strconv.Itoa(userData.ID), Name: userData.Name};
    // if !exists {
    //     return nil, grpc.Errorf(grpc.Code(nil), "user not found")
    // }
    
    return user, nil
}

func (s *Server) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    // log.Printf("Received CreateUser request: %s, %s", req.Name, req.Email)
    
    user := &UserResponse{
        UserId: "user",
        Name:   req.Name,
        Email:  req.Email,
    }
    
    s.users[user.UserId] = user
    return user, nil
}

// func generateID() string {
//     return "user_";
// }

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