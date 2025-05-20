package service

import (
	pb "Taurus/internal/controller/gRPC/proto/order"
	"Taurus/pkg/grpc/server"
	"context"
	"log"

	"google.golang.org/grpc"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
}

func NewOrderService() *OrderService {
	log.Printf("创建 OrderService 实例")
	return &OrderService{}
}

// QueryOrders 实现查询订单列表
func (s *OrderService) QueryOrders(ctx context.Context, req *pb.QueryOrdersRequest) (*pb.QueryOrdersResponse, error) {
	// TODO: 实现查询逻辑
	return &pb.QueryOrdersResponse{
		Orders: []*pb.OrderDetail{
			{
				OrderId:    "123",
				UserId:     "user1",
				Amount:     100.0,
				Status:     "paid",
				CreateTime: "2024-01-01",
				UpdateTime: "2024-01-01",
				Items: []*pb.OrderItem{
					{
						ItemId:   "item1",
						ItemName: "商品1",
						Quantity: 1,
						Price:    100.0,
					},
				},
			},
		},
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetOrderDetail 实现获取订单详情
func (s *OrderService) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailRequest) (*pb.GetOrderDetailResponse, error) {
	// TODO: 实现查询逻辑
	return &pb.GetOrderDetailResponse{
		Order: &pb.Order{
			OrderId:    req.OrderId,
			UserId:     "user1",
			Amount:     100.0,
			Status:     "paid",
			CreateTime: "2024-01-01",
			UpdateTime: "2024-01-01",
			Items: []*pb.OrderItem{
				{
					ItemId:   "item1",
					ItemName: "商品1",
					Quantity: 1,
					Price:    100.0,
				},
			},
		},
	}, nil
}

// RegisterService 实现服务注册接口
func (s *OrderService) RegisterService(server *grpc.Server) {
	log.Printf("注册 OrderService 到 gRPC 服务器")
	pb.RegisterOrderServiceServer(server, s)
}

func init() {
	server.RegisterService("order", NewOrderService())
}
