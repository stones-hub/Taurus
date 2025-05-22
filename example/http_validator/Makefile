.PHONY: build run test clean deps fmt lint help

# 默认目标
all: build

# 构建应用
build:
	@echo "构建 HTTP 验证示例..."
	@go build -o bin/http_validator main.go

# 运行应用
run: build
	@echo "启动 HTTP 验证示例..."
	@./bin/http_validator

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 清理构建产物
clean:
	@echo "清理构建产物..."
	@rm -rf bin/

# 安装依赖
deps:
	@echo "安装依赖..."
	@go mod tidy

# 代码格式化
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 代码检查
lint:
	@echo "检查代码..."
	@go vet ./...

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  make build    - 构建应用"
	@echo "  make run      - 构建并运行应用"
	@echo "  make test     - 运行测试"
	@echo "  make clean    - 清理构建产物"
	@echo "  make deps     - 安装依赖"
	@echo "  make fmt      - 格式化代码"
	@echo "  make lint     - 检查代码"
	@echo "  make help     - 显示帮助信息" 