# Telemetry 调用链监控示例

本目录包含了使用 `pkg/telemetry` 包进行调用链追踪的示例代码。

## 依赖安装

```bash
# OpenTelemetry 基础依赖
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/sdk

# 数据库相关
go get github.com/go-sql-driver/mysql
go get github.com/go-redis/redis/v8
```

## 示例说明

### 1. HTTP 服务追踪
文件：`http/main.go`

展示了如何：
- 初始化 telemetry provider
- 在 HTTP 处理器中创建 span
- 追踪请求处理过程

测试：
```bash
# 启动服务
go run http/main.go

# 发送请求
curl http://localhost:8080/hello
```

### 2. gRPC 服务追踪
文件：`grpc/main.go`

展示了如何：
- 在 gRPC 服务中使用追踪
- 创建服务级别的 tracer
- 追踪 RPC 调用

### 3. Redis 操作追踪
文件：`redis/main.go`

展示了如何：
- 追踪 Redis 操作
- 使用 span 记录缓存操作
- 传递上下文

运行前确保 Redis 已启动：
```bash
docker run -d --name redis -p 6379:6379 redis
```

### 4. MySQL 操作追踪
文件：`mysql/main.go`

展示了如何：
- 追踪数据库操作
- 在查询中使用上下文
- 记录查询结果

运行前准备数据库：
```bash
docker run -d --name mysql \
    -e MYSQL_ROOT_PASSWORD=password \
    -e MYSQL_DATABASE=test \
    -p 3306:3306 \
    mysql:8

# 创建测试表
mysql -h 127.0.0.1 -u root -ppassword test <<EOF
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
INSERT INTO users (id, name) VALUES (1, 'Test User');
EOF
```

## 查看追踪结果

所有示例都会将追踪数据发送到 Jaeger。启动 Jaeger：

```bash
docker run -d --name jaeger \
    -e COLLECTOR_OTLP_ENABLED=true \
    -p 16686:16686 \
    -p 4317:4317 \
    jaegertracing/all-in-one:latest
```

访问 Jaeger UI：http://localhost:16686

## 注意事项

1. 这些示例都是独立的，展示了不同场景下的追踪使用
2. 每个示例都使用了最简单的方式来展示追踪功能
3. 生产环境中需要添加更多的错误处理和配置选项
4. 示例中省略了一些生产环境必需的设置（如连接池配置等） 