package test

import (
	"context"
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

	// 创建用户服务客户端
	userClient := userpb.NewUserServiceClient(conn)

	// 测试获取用户信息
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
}

//nolint:unused // 保留此函数供后续使用
func setToken(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": token,
	})
	return metadata.NewOutgoingContext(ctx, md)
}
