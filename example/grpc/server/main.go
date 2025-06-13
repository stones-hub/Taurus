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
	"Taurus/pkg/grpc/server"
	"context"
	"fmt"
	"log"

	pb "Taurus/example/grpc/proto/user"
)

// UserServer implements the user service defined in the proto file.
// It provides methods for user management operations.
type UserServer struct {
	pb.UnimplementedUserServiceServer
}

// GetUser retrieves user information by ID.
// This is a simple implementation that returns mock data.
//
// Parameters:
//   - ctx: Context for the request
//   - req: GetUserRequest containing the user ID
//
// Returns:
//   - *pb.GetUserResponse: User information
//   - error: Any error that occurred during the operation
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Mock user data retrieval from database
	user := &pb.GetUserResponse{
		Id:    req.Id,
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
	}
	return user, nil
}

// CreateUser creates a new user with the provided information.
// This is a simple implementation that returns mock data.
//
// Parameters:
//   - ctx: Context for the request
//   - req: CreateUserRequest containing the user information
//
// Returns:
//   - *pb.CreateUserResponse: Created user information
//   - error: Any error that occurred during the operation
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Mock user creation
	user := &pb.CreateUserResponse{
		Id:    1,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	return user, nil
}

// main is the entry point of the application.
// It sets up and starts the gRPC server with the user service.
func main() {
	// Create a new gRPC server instance
	srv, cleanup, err := server.NewServer(
		server.WithAddress(":50051"),
	)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	defer cleanup()

	// Register the user service with the server
	pb.RegisterUserServiceServer(srv.Server(), &UserServer{})

	// Start the server
	fmt.Println("Starting gRPC server on :50051")
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
