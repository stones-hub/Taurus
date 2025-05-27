package main

import (
	"context"
	"log"

	"Taurus/pkg/telemetry"

	"github.com/go-redis/redis/v8"
)

func main() {
	// 初始化 provider
	provider, err := telemetry.NewOTelProvider(
		telemetry.WithServiceName("redis-demo"),
		telemetry.WithEnvironment("dev"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer provider.Shutdown(context.Background())

	// 创建 Redis 客户端
	opts := &redis.Options{
		Addr: "localhost:6379",
	}
	client := telemetry.WrapRedis(opts)

	// 测试 Redis 连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	// 设置一个值
	if err := client.Set(ctx, "hello", "world", 0).Err(); err != nil {
		log.Fatal(err)
	}

	// 获取值
	val, err := client.Get(ctx, "hello").Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Value from Redis: %s", val)
}
