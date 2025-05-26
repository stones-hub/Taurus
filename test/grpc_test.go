package test

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	ordpb "Taurus/internal/controller/gRPC/proto/order"
	userpb "Taurus/internal/controller/gRPC/proto/user"
	"Taurus/pkg/grpc/client"
	"Taurus/pkg/grpc/client/interceptor"

	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

func TestOrderService(t *testing.T) {
	// 创建客户端，添加完整配置
	c, err := client.NewClient(
		// 基础配置
		client.WithTimeout(5*time.Second),
		client.WithInsecure(), // 测试环境使用非安全连接

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
		client.WithUnaryInterceptor(interceptor.AuthInterceptor("Bearer 123456")),
		client.WithStreamInterceptor(interceptor.StreamAuthInterceptor("Bearer 123456")),
	)
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}
	defer c.Close()

	// 获取连接（非流式）
	conn, err := c.GetConn("localhost:50051", false)
	if err != nil {
		t.Fatalf("获取连接失败: %v", err)
	}
	defer c.ReleaseConn(conn)

	// 创建服务客户端
	orderClient := ordpb.NewOrderServiceClient(conn)

	// 测试查询订单列表
	t.Run("QueryOrders", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := orderClient.QueryOrders(ctx, &ordpb.QueryOrdersRequest{
			StartDate: "2024-01-01",
			EndDate:   "2024-01-31",
			Page:      1,
			PageSize:  10,
		})
		if err != nil {
			t.Fatalf("查询订单失败: %v", err)
		}
		t.Logf("查询结果: %+v", resp)

		// 测试无效请求
		invalidReq := &ordpb.QueryOrdersRequest{
			StartDate: "2024-01-01 00:00:00",
			EndDate:   "2024-01-31 23:59:59",
			Page:      1,
			PageSize:  -1,
		}
		resp, err = orderClient.QueryOrders(ctx, invalidReq)
		if err != nil {
			t.Logf("预期的错误: %v", err)
		}
		t.Logf("无效请求结果: %+v", resp)
	})

	// 测试获取订单详情
	t.Run("GetOrderDetail", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := orderClient.GetOrderDetail(ctx, &ordpb.GetOrderDetailRequest{
			OrderId: "123",
		})
		if err != nil {
			t.Fatalf("获取订单详情失败: %v", err)
		}
		t.Logf("订单详情: %+v", resp)
	})
}

func TestUserService(t *testing.T) {
	// 创建客户端，添加完整配置
	c, err := client.NewClient(
		// 基础配置
		client.WithTimeout(5*time.Second),
		client.WithInsecure(), // 测试环境使用非安全连接

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
			Time:                10 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),

		// 添加认证拦截器
		client.WithUnaryInterceptor(interceptor.AuthInterceptor("Bearer 123456")),
		client.WithStreamInterceptor(interceptor.StreamAuthInterceptor("Bearer 123456")),
	)
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}
	defer c.Close()

	// 获取连接（非流式）
	conn, err := c.GetConn("localhost:50051", false)
	if err != nil {
		t.Fatalf("获取连接失败: %v", err)
	}
	defer c.ReleaseConn(conn)

	streamConn, err := c.GetConn("localhost:50051", true)
	if err != nil {
		t.Fatalf("获取流式连接失败: %v", err)
	}
	defer c.ReleaseConn(streamConn)

	// 创建用户服务客户端
	userClient := userpb.NewUserServiceClient(conn)
	streamUserClient := userpb.NewUserServiceClient(streamConn)

	// 测试一元调用 - 获取用户信息
	t.Run("GetUserInfo", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := userClient.GetUserInfo(ctx, &userpb.GetUserInfoRequest{
			UserId: 1,
		})
		if err != nil {
			t.Fatalf("获取用户信息失败: %v", err)
		}
		t.Logf("用户信息: %+v", resp)
	})

	// 测试服务端流式调用 - 获取用户列表
	t.Run("GetUserList", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stream, err := streamUserClient.GetUserList(ctx, &userpb.GetUserListRequest{
			UserIds:  []int64{1, 2, 3, 4, 5},
			PageSize: 10,
			PageNum:  1,
		})
		if err != nil {
			t.Fatalf("获取用户列表流失败: %v", err)
		}

		for {
			user, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("接收用户数据失败: %v", err)
			}
			t.Logf("收到用户数据: %+v", user)
		}
	})

	// 测试客户端流式调用 - 批量创建用户
	t.Run("BatchCreateUsers", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stream, err := streamUserClient.BatchCreateUsers(ctx)
		if err != nil {
			t.Fatalf("创建批量用户流失败: %v", err)
		}

		// 发送多个用户创建请求
		users := []struct {
			username string
			email    string
			age      int32
			password string
		}{
			{"user1", "user1@example.com", 25, "password1"},
			{"user2", "user2@example.com", 30, "password2"},
			{"user3", "user3@example.com", 35, "password3"},
		}

		for _, u := range users {
			if err := stream.Send(&userpb.CreateUserRequest{
				Username: u.username,
				Email:    u.email,
				Age:      u.age,
				Password: u.password,
			}); err != nil {
				t.Fatalf("发送用户创建请求失败: %v", err)
			}
		}

		// 关闭发送并获取响应
		resp, err := stream.CloseAndRecv()
		if err != nil {
			t.Fatalf("接收批量创建响应失败: %v", err)
		}
		t.Logf("批量创建结果: 成功=%d, 失败=%d", resp.SuccessCount, resp.FailedCount)
		for _, msg := range resp.ErrorMessages {
			t.Logf("错误信息: %s", msg)
		}
	})

	// 测试双向流式调用 - 用户信息同步
	t.Run("SyncUserInfo", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stream, err := streamUserClient.SyncUserInfo(ctx)
		if err != nil {
			t.Fatalf("创建同步流失败: %v", err)
		}

		// 创建错误通道
		errChan := make(chan error)
		// 创建完成通道
		doneChan := make(chan bool)

		// 启动接收协程
		go func() {
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					doneChan <- true
					return
				}
				if err != nil {
					errChan <- fmt.Errorf("接收同步响应失败: %v", err)
					return
				}
				t.Logf("收到同步响应: %+v", resp)
			}
		}()

		// 发送同步请求
		syncRequests := []struct {
			userId    int64
			username  string
			email     string
			age       int32
			operation string
		}{
			{1, "user1_updated", "user1@example.com", 26, "update"},
			{2, "user2_updated", "user2@example.com", 31, "update"},
			{3, "user3", "user3@example.com", 35, "delete"},
		}

		for _, req := range syncRequests {
			if err := stream.Send(&userpb.UserInfoSync{
				UserId:    req.userId,
				Username:  req.username,
				Email:     req.email,
				Age:       req.age,
				Operation: req.operation,
			}); err != nil {
				t.Fatalf("发送同步请求失败: %v", err)
			}
			time.Sleep(100 * time.Millisecond) // 模拟间隔
		}

		// 关闭发送
		if err := stream.CloseSend(); err != nil {
			t.Fatalf("关闭发送失败: %v", err)
		}

		// 等待完成或错误
		select {
		case <-doneChan:
			t.Log("同步完成")
		case err := <-errChan:
			t.Fatalf("同步过程出错: %v", err)
		case <-ctx.Done():
			t.Fatalf("同步超时")
		}
	})
}

//nolint:unused // 保留此函数供后续使用
func setToken(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": token,
	})
	return metadata.NewOutgoingContext(ctx, md)
}
