package main

import (
	"context"
	"log"
	"time"

	pb "Taurus/example/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := pb.NewCalculatorClient(conn)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用加法服务
	addResp, err := client.Add(ctx, &pb.AddRequest{A: 10, B: 20})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Printf("Add result: %d", addResp.Result)

	// 调用减法服务
	subtractResp, err := client.Subtract(ctx, &pb.SubtractRequest{A: 30, B: 15})
	if err != nil {
		log.Fatalf("could not subtract: %v", err)
	}
	log.Printf("Subtract result: %d", subtractResp.Result)
}
