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
	"context"
	"log"

	"Taurus/pkg/redisx"
	"Taurus/pkg/telemetry"
)

func main() {
	// 1. 初始化追踪器提供者
	provider, _, err := telemetry.NewOTelProvider(
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
