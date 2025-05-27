package main

import (
	"context"
	"log"

	"Taurus/pkg/redisx"
	"Taurus/pkg/telemetry"
)

func main() {
	// 1. 初始化追踪器提供者
	provider, err := telemetry.NewOTelProvider(
		telemetry.WithServiceName("redis-demo"),
		telemetry.WithServiceVersion("v0.1.0"),
		telemetry.WithEnvironment("dev"),
	)
	if err != nil {
		log.Fatalf("init telemetry provider failed: %v", err)
	}
	defer provider.Shutdown(context.Background())

	// 2. 初始化 Redis
	redisTracer := provider.Tracer("redis-client")
	redisClient := redisx.InitRedis(redisx.RedisConfig{
		Addrs:    []string{"localhost:6379"},
		Password: "",
		DB:       0,
	})

	// 添加追踪 Hook
	redisClient.AddHook(&telemetry.RedisHook{
		Tracer: redisTracer,
	})

	// 3. 执行一些 Redis 操作
	ctx := context.Background()
	if err := redisClient.Set(ctx, "test_key", "test_value", 0); err != nil {
		log.Printf("set key failed: %v", err)
	}

	value, err := redisClient.Get(ctx, "test_key")
	if err != nil {
		log.Printf("get key failed: %v", err)
	}
	log.Printf("get value: %v", value)

	log.Printf("Redis demo completed")
}
