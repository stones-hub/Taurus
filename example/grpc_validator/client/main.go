// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

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
		// 基础配置
		client.WithTimeout(5*time.Second),
		client.WithInsecure(),

		// 连接池配置
		client.WithPoolConfig(
			5,              // maxIdle
			50,             // maxOpen
			30*time.Minute, // maxLifetime
			10*time.Minute, // maxIdleTime
			1000,           // maxLoad
		),
	)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer c.Close()

	// 获取连接（非流式）
	conn, err := c.GetConn("localhost:50051", false)
	if err != nil {
		log.Fatalf("获取连接失败: %v", err)
	}
	defer c.ReleaseConn(conn)

	// 创建用户服务客户端
	userClient := pb.NewUserServiceClient(conn)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试创建用户 - 有效请求
	createUserResp, err := userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:     "张三",
		Email:    "zhangsan@example.com",
		Age:      25,
		Password: "123456",
	})
	if err != nil {
		log.Printf("创建用户失败: %v", err)
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
		log.Printf("预期的无效请求错误: %v", err)
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
		log.Printf("更新用户失败: %v", err)
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
		log.Printf("预期的无效请求错误: %v", err)
	} else {
		fmt.Printf("UpdateUser response: %+v\n", updateUserResp)
	}
}
