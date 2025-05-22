package main

import (
	"Taurus/pkg/grpc/client"
	"context"
	"fmt"
	"log"
	"time"

	pb "Taurus/example/grpc/proto/user"

	"google.golang.org/grpc/keepalive"
)

func main() {
	// 创建gRPC客户端
	c, err := client.NewClient(
		client.WithAddress("localhost:50051"),
		client.WithTimeout(5*time.Second),
		client.WithInsecure(),
		client.WithToken("123456"),
		client.WithKeepAlive(&keepalive.ClientParameters{
			Time:                10 * time.Second, // 发送 keepalive 的时间间隔
			Timeout:             5 * time.Second,  // keepalive 超时时间
			PermitWithoutStream: true,             // 允许在没有活跃流的情况下发送 keepalive
		}),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	// 创建用户服务客户端
	userClient := pb.NewUserServiceClient(c.Conn())

	// 创建上下文
	ctx := context.Background()

	// 调用GetUser方法
	getUserResp, err := userClient.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	fmt.Printf("GetUser response: %+v\n", getUserResp)

	// 调用CreateUser方法
	createUserResp, err := userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "李四",
		Email: "lisi@example.com",
		Age:   30,
	})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	fmt.Printf("CreateUser response: %+v\n", createUserResp)
}
