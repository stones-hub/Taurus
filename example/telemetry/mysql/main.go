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

	"Taurus/pkg/db"
	"Taurus/pkg/telemetry"
)

// User 用户模型
type User struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"type:varchar(32)"`
	Age  int    `gorm:"type:int"`
}

func main() {
	// 1. 初始化追踪器提供者
	provider, _, err := telemetry.NewOTelProvider(
		telemetry.WithServiceName("mysql-demo"),
		telemetry.WithServiceVersion("v0.1.0"),
		telemetry.WithEnvironment("dev"),
	)
	if err != nil {
		log.Fatalf("init telemetry provider failed: %v", err)
	}
	defer provider.Shutdown(context.Background())

	// 2. 初始化 MySQL
	mysqlTracer := provider.Tracer("mysql-client")
	db.InitDB("default", "mysql", "root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		db.NewDBCustomLogger(db.DBLoggerConfig{
			LogLevel: 4,
		}), 3, 5)
	defer db.CloseDB()

	// 为默认数据库添加追踪
	for _, db := range db.DbList() {
		if err := db.Use(&telemetry.GormTracingHook{
			Tracer: mysqlTracer,
		}); err != nil {
			log.Fatalf("use tracing hook failed: %v", err)
		}
	}

	// 3. 执行一些数据库操作
	var user User
	if err := db.Find("default", &user, 1); err != nil {
		log.Printf("query user failed: %v", err)
	}

	log.Printf("MySQL demo completed")
}
