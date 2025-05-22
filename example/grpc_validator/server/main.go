package main

import (
	"Taurus/pkg/grpc/server"
	"Taurus/pkg/grpc/server/interceptor"
	"context"
	"fmt"
	"log"

	pb "Taurus/example/grpc_validator/proto/user"
)

// UserServer 实现用户服务
type UserServer struct {
	pb.UnimplementedUserServiceServer
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

// UpdateUser 更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// 模拟更新用户
	user := &pb.UpdateUserResponse{
		Id:    req.Id,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	return user, nil
}

func main() {
	// 创建gRPC服务器，添加验证拦截器
	srv, cleanup, err := server.NewServer(
		server.WithAddress(":50051"),
		server.WithUnaryInterceptor(
			interceptor.UnaryServerValidationInterceptor(),
		),
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
