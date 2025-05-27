# Telemetry 分布式追踪组件

基于 OpenTelemetry 的分布式追踪组件，提供了简单易用的 API 来实现应用程序的链路追踪功能。

## 一、使用指南

### 1. 安装依赖

```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/sdk
```

### 2. 初始化追踪提供者

```go
import "your-project/pkg/telemetry"

provider, err := telemetry.NewOTelProvider(
    telemetry.WithServiceName("your-service"),    // 设置服务名称
    telemetry.WithServiceVersion("1.0.0"),        // 设置服务版本
    telemetry.WithEnvironment("production"),       // 设置环境
    telemetry.WithEndpoint("localhost:4317"),      // 设置Collector地址
    telemetry.WithExportProtocol("grpc"),         // 设置协议（支持grpc/http/json）
)
if err != nil {
    log.Fatal(err)
}
defer provider.Shutdown(context.Background())
```

### 3. 创建追踪

```go
// 创建追踪器
tracer := provider.Tracer("component-name")

// 创建span
ctx, span := tracer.Start(context.Background(), "operation-name")
defer span.End()

// 你的业务代码
```

### 4. MySQL和Redis自动追踪

```go
// MySQL追踪
db, err := telemetry.WrapMySQL(rawDB, "mysql-service-name")

// Redis追踪
rdb, err := telemetry.WrapRedis(rawClient, "redis-service-name")
```

## 二、单机部署方案

### 1. 组件架构

单机部署只需要以下核心组件：
- OpenTelemetry Collector：接收和处理追踪数据
- Jaeger All-in-One：包含 Collector、Query UI 和存储功能

### 2. 部署配置

```yaml
# docker-compose.yml
version: "3"
services:
  # 1. Jaeger All-in-One
  jaeger:
    image: jaegertracing/all-in-one:latest
    environment:
      - COLLECTOR_OTLP_ENABLED=true
      - BADGER_EPHEMERAL=false
      - SPAN_STORAGE_TYPE=badger
      - BADGER_DIRECTORY_VALUE=/badger/data
      - BADGER_DIRECTORY_KEY=/badger/key
    ports:
      - "16686:16686"  # UI
      - "14250:14250"  # gRPC
      - "14268:14268"  # HTTP
      - "6831:6831/udp"  # jaeger.thrift 
      - "6832:6832/udp"  # jaeger.thrift 
      - "4317:4317"   # OTLP gRPC
      - "4318:4318"   # OTLP HTTP
    volumes:
      - jaeger-data:/badger

  # 2. OpenTelemetry Collector（可选，如果需要更多功能）
  otel-collector:
    image: otel/opentelemetry-collector:latest
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "8888:8888"   # 监控指标

volumes:
  jaeger-data:
    driver: local
```

### 3. Collector 配置（可选）

```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 1s
    send_batch_size: 1024

exporters:
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [jaeger]
```

### 4. 使用说明

1. 启动服务：
```bash
docker-compose up -d
```

2. 访问 Jaeger UI：
- 地址：`http://localhost:16686`
- 功能：
  - 查看调用链路图
  - 分析服务依赖
  - 查看调用时间
  - 查询和过滤追踪数据

3. 应用程序配置：
```go
provider, err := telemetry.NewOTelProvider(
    telemetry.WithServiceName("your-service"),
    telemetry.WithEndpoint("localhost:4317"),
)
```

### 5. 数据管理

1. 存储配置：
- 默认使用 Badger 存储
- 数据保留 72 小时
- 数据存储在 Docker 卷中

2. 调整保留时间：
```yaml
environment:
  - BADGER_RETENTION=168h  # 保留7天
```

### 6. 注意事项

1. 资源建议：
- 建议 2 CPU，4GB 内存以上
- 存储空间 20GB 以上
- 适合单机每天百万级 spans 规模

2. 数据备份：
- 定期备份 badger 数据目录
- 备份配置文件

3. 扩展性：
- 如果数据量增长较大，可以平滑迁移到分布式方案
- 支持后续扩展为高可用架构 