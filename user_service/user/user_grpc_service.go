package user

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
    UnimplementedUserServiceServer
    users map[string]*UserResponse
}

func NewServer() *Server {
    return &Server{
        users: make(map[string]*UserResponse),
    }
}

func (s *Server) GetUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
    // log.Printf("Received GetUser request for user_id: %s", req.UserId)
    
    // user, exists := s.users[req.UserId]
    user := &UserResponse {UserId: "12"};
    // if !exists {
    //     return nil, grpc.Errorf(grpc.Code(nil), "user not found")
    // }
    
    
    return user, nil
}

func (s *Server) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    log.Printf("Received CreateUser request: %s, %s", req.Name, req.Email)
    
    user := &UserResponse{
        UserId: generateID(),
        Name:   req.Name,
        Email:  req.Email,
    }
    
    s.users[user.UserId] = user
    return user, nil
}

func generateID() string {
    return "user_";
}

func StartUserServer() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    s := grpc.NewServer()
    RegisterUserServiceServer(s, NewServer())
    
    log.Printf("UserService server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}