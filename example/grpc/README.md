# gRPC 示例

这是一个简单的 gRPC 示例，实现了一个计算器服务，支持加法和减法运算。

## 前置要求

1. 安装 Protocol Buffers 编译器：
```bash
# macOS
brew install protobuf

# Ubuntu
apt-get install protobuf-compiler
```

2. 安装 Go 的 protobuf 插件：
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 使用说明

1. 生成 protobuf 代码：
```bash
make proto
```

2. 启动服务器：
```bash
make run-server
```

3. 在另一个终端中运行客户端：
```bash
make run-client
```

## 代码结构

- `proto/calculator.proto`: 服务定义文件
- `server/main.go`: 服务端实现
- `client/main.go`: 客户端实现
- `Makefile`: 构建和运行命令

## 示例说明

这个示例实现了一个简单的计算器服务，提供以下功能：
- 加法运算 (Add)
- 减法运算 (Subtract)

服务端监听在 `localhost:50051`，客户端会连接到这个地址并发送请求。

## 扩展建议

1. 添加更多运算方法（乘法、除法等）
2. 添加错误处理
3. 实现流式 RPC
4. 添加 TLS 支持
5. 添加服务发现机制 