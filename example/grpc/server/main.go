package main

import (
	"context"
	"log"
	"net"

	pb "Taurus/example/grpc/proto"

	"google.golang.org/grpc"
)

// CalculatorServer 实现 Calculator 服务
type CalculatorServer struct {
	pb.UnimplementedCalculatorServer
}

// Add 实现加法运算
func (s *CalculatorServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	result := req.A + req.B
	return &pb.AddResponse{Result: result}, nil
}

// Subtract 实现减法运算
func (s *CalculatorServer) Subtract(ctx context.Context, req *pb.SubtractRequest) (*pb.SubtractResponse, error) {
	result := req.A - req.B
	return &pb.SubtractResponse{Result: result}, nil
}

func main() {
	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册服务
	pb.RegisterCalculatorServer(server, &CalculatorServer{})

	// 启动服务器
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Server is running on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
