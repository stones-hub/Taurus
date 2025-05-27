package main

import (
	"context"
	"log"

	"Taurus/pkg/telemetry"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 初始化 provider
	provider, err := telemetry.NewOTelProvider(
		telemetry.WithServiceName("mysql-demo"),
		telemetry.WithEnvironment("dev"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer provider.Shutdown(context.Background())

	// 创建带追踪的数据库连接
	dsn := "root:password@tcp(localhost:3306)/test?parseTime=true"
	db, err := telemetry.WrapMySQL(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建带追踪的数据库操作包装器
	tracedDB := telemetry.NewTracedDB(db)

	// 测试数据库连接
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	// 执行查询
	var name string
	rows, err := tracedDB.QueryContext(ctx,
		"SELECT name FROM users WHERE id = ?",
		1,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		log.Printf("User name: %s", name)
	} else {
		log.Println("User not found")
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
