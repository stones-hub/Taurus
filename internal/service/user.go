package service

import (
	pb "Taurus/internal/controller/gRPC/user"
	"Taurus/pkg/grpc/server"
	"context"

	"google.golang.org/grpc"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

// GetUserInfo 实现获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	// TODO: 实现查询逻辑
	return &pb.GetUserInfoResponse{
		UserId:   req.UserId,
		Username: "test_user",
		Email:    "test@example.com",
		Age:      25,
	}, nil
}

func (s *UserService) Register(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, s)
}

func init() {
	server.RegisterService("user", NewUserService())
}
