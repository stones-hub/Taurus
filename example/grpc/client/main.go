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
	"Taurus/pkg/grpc/client/interceptor"
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

		// 保活配置
		client.WithKeepAlive(&keepalive.ClientParameters{
			Time:                10 * time.Second, // 发送 keepalive 的时间间隔
			Timeout:             5 * time.Second,  // keepalive 超时时间
			PermitWithoutStream: true,             // 允许在没有活跃流的情况下发送 keepalive
		}),

		// 添加认证拦截器
		client.WithUnaryInterceptor(interceptor.AuthInterceptor("123456")),
		client.WithStreamInterceptor(interceptor.StreamAuthInterceptor("123456")),
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

	// 调用GetUser方法
	getUserResp, err := userClient.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	if err != nil {
		log.Fatalf("获取用户失败: %v", err)
	}
	fmt.Printf("GetUser response: %+v\n", getUserResp)

	// 调用CreateUser方法
	createUserResp, err := userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "李四",
		Email: "lisi@example.com",
		Age:   30,
	})
	if err != nil {
		log.Fatalf("创建用户失败: %v", err)
	}
	fmt.Printf("CreateUser response: %+v\n", createUserResp)
}
