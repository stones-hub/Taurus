package main

import (
	"Taurus/pkg/grpc/server"
	"context"
	"fmt"
	"log"

	pb "Taurus/example/grpc/proto/user"
)

// UserServer 实现用户服务
type UserServer struct {
	pb.UnimplementedUserServiceServer
}

// GetUser 获取用户信息
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// 模拟从数据库获取用户
	user := &pb.GetUserResponse{
		Id:    req.Id,
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
	}
	return user, nil
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// 模拟创建用户
	user := &pb.CreateUserResponse{
		Id:    1,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	return user, nil
}

func main() {
	// 创建gRPC服务器
	srv, cleanup, err := server.NewServer(
		server.WithAddress(":50051"),
	)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	defer cleanup()

	// 注册服务
	pb.RegisterUserServiceServer(srv.Server(), &UserServer{})

	// 启动服务器
	fmt.Println("Starting gRPC server on :50051")
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
