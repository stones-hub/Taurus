package service

import (
	"Taurus/pkg/grpc/server"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "Taurus/internal/controller/gRPC/proto/user"
)

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	log.Printf("创建 UserService 实例")
	return &UserService{}
}

// GetUserInfo 一元调用 - 获取单个用户信息
func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	log.Printf("收到 GetUserInfo 请求: user_id=%d", req.UserId)

	// TODO: 实现实际的数据库查询逻辑
	return &pb.GetUserInfoResponse{
		UserId:    req.UserId,
		Username:  "test_user",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// GetUserList 服务端流式调用 - 批量获取用户信息
func (s *UserService) GetUserList(req *pb.GetUserListRequest, stream pb.UserService_GetUserListServer) error {
	log.Printf("收到 GetUserList 请求: page_size=%d, page_num=%d", req.PageSize, req.PageNum)

	// 参数验证
	if req.PageSize <= 0 || req.PageNum <= 0 {
		return status.Error(codes.InvalidArgument, "页码和分页大小必须大于0")
	}

	// 模拟分页返回用户数据
	for i := 0; i < int(req.PageSize); i++ {
		// 检查上下文是否已取消
		if err := stream.Context().Err(); err != nil {
			return status.Error(codes.Canceled, "客户端取消了请求")
		}

		// 构造用户数据
		user := &pb.GetUserInfoResponse{
			UserId:    int64(i + 1),
			Username:  fmt.Sprintf("user_%d", i+1),
			Email:     fmt.Sprintf("user%d@example.com", i+1),
			Age:       20 + int32(i%10),
			CreatedAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		// 发送数据
		if err := stream.Send(user); err != nil {
			return status.Errorf(codes.Internal, "发送数据失败: %v", err)
		}

		// 模拟处理延迟
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// BatchCreateUsers 客户端流式调用 - 批量创建用户
func (s *UserService) BatchCreateUsers(stream pb.UserService_BatchCreateUsersServer) error {
	log.Printf("开始处理 BatchCreateUsers 请求")

	var (
		successCount int32
		failedCount  int32
		users        []*pb.GetUserInfoResponse
		errors       []string
	)

	// 接收并处理客户端流数据
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// 客户端已完成发送
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "接收数据失败: %v", err)
		}

		// 验证数据
		if err := validateCreateUserRequest(req); err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("用户 %s 验证失败: %v", req.Username, err))
			continue
		}

		// 模拟创建用户
		user := &pb.GetUserInfoResponse{
			UserId:    int64(len(users) + 1),
			Username:  req.Username,
			Email:     req.Email,
			Age:       req.Age,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		users = append(users, user)
		successCount++
	}

	// 返回批处理结果
	return stream.SendAndClose(&pb.BatchCreateUsersResponse{
		Users:         users,
		SuccessCount:  successCount,
		FailedCount:   failedCount,
		ErrorMessages: errors,
	})
}

// SyncUserInfo 双向流式调用 - 实时用户信息同步
func (s *UserService) SyncUserInfo(stream pb.UserService_SyncUserInfoServer) error {
	log.Printf("开始处理 SyncUserInfo 请求")

	// 创建错误通道
	errChan := make(chan error)

	// 启动接收协程
	go func() {
		for {
			// 接收客户端消息
			req, err := stream.Recv()
			if err == io.EOF {
				errChan <- nil
				return
			}
			if err != nil {
				errChan <- status.Errorf(codes.Internal, "接收数据失败: %v", err)
				return
			}

			// 处理同步请求
			resp := &pb.UserInfoSync{
				UserId:    req.UserId,
				Username:  req.Username,
				Email:     req.Email,
				Age:       req.Age,
				Timestamp: time.Now().Unix(),
				Operation: req.Operation,
			}

			// 发送响应
			if err := stream.Send(resp); err != nil {
				errChan <- status.Errorf(codes.Internal, "发送数据失败: %v", err)
				return
			}
		}
	}()

	// 等待处理完成或出错
	err := <-errChan
	return err
}

// validateCreateUserRequest 验证创建用户请求
func validateCreateUserRequest(req *pb.CreateUserRequest) error {
	if req.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if req.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	if req.Age <= 0 || req.Age > 120 {
		return fmt.Errorf("年龄必须在1-120之间")
	}
	if req.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	return nil
}

func (s *UserService) RegisterService(server *grpc.Server) {
	log.Printf("注册 UserService 到 gRPC 服务器")
	pb.RegisterUserServiceServer(server, s)
}

func init() {
	server.RegisterService("user", NewUserService())
}

/*
// 生成 proto 文件
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/controller/gRPC/proto/user/user.proto
*/
