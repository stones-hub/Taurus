package service

import (
	pb "Taurus/internal/controller/gRPC/proto/user"
	"Taurus/pkg/grpc/server"
	"context"
	"log"

	"google.golang.org/grpc"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	log.Printf("创建 UserService 实例")
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

func (s *UserService) RegisterService(server *grpc.Server) {
	log.Printf("注册 UserService 到 gRPC 服务器")
	pb.RegisterUserServiceServer(server, s)
}

func init() {
	server.RegisterService("user", NewUserService())
}
