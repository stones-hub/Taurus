package main

import (
	"Taurus/pkg/grpc/client"
	"context"
	"fmt"
	"log"
	"time"

	pb "Taurus/example/grpc_validator/proto/user"
)

func main() {
	// 创建gRPC客户端
	c, err := client.NewClient(
		client.WithAddress("localhost:50051"),
		client.WithTimeout(5*time.Second),
		client.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	// 创建用户服务客户端
	userClient := pb.NewUserServiceClient(c.Conn())

	// 创建上下文
	ctx := context.Background()

	// 测试创建用户 - 有效请求
	createUserResp, err := userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:     "张三",
		Email:    "zhangsan@example.com",
		Age:      25,
		Password: "123456",
	})
	if err != nil {
		log.Printf("failed to create user: %v", err)
	} else {
		fmt.Printf("CreateUser response: %+v\n", createUserResp)
	}

	// 测试创建用户 - 无效请求（年龄小于0）
	createUserResp, err = userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:     "李四",
		Email:    "invalid-email",
		Age:      -1,
		Password: "123",
	})
	if err != nil {
		log.Printf("expected error for invalid request: %v", err)
	} else {
		fmt.Printf("CreateUser response: %+v\n", createUserResp)
	}

	// 测试更新用户 - 有效请求
	updateUserResp, err := userClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:    1,
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   26,
	})
	if err != nil {
		log.Printf("failed to update user: %v", err)
	} else {
		fmt.Printf("UpdateUser response: %+v\n", updateUserResp)
	}

	// 测试更新用户 - 无效请求（ID为0）
	updateUserResp, err = userClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:    0,
		Name:  "李四",
		Email: "invalid-email",
		Age:   -1,
	})
	if err != nil {
		log.Printf("expected error for invalid request: %v", err)
	} else {
		fmt.Printf("UpdateUser response: %+v\n", updateUserResp)
	}
}
