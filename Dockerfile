# 使用官方的 Go 语言镜像作为基础镜像
FROM golang:1.24-alpine AS builder

# 设置工作目录, 容器启动后会进入该目录
ARG WORKDIR
WORKDIR ${WORKDIR}

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 将项目的所有文件复制到工作目录
COPY . .

# 编译 Go 应用程序
RUN go build -o main cmd/main.go

# 使用一个更小的基础镜像来运行应用程序
FROM alpine:latest

# 安装必要的依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
ARG WORKDIR
WORKDIR ${WORKDIR}

# 从构建阶段复制编译好的二进制文件
# 在 Dockerfile 中，COPY 指令用于将文件或目录从构建阶段复制到新镜像中。
# 这里，--from=builder 指定从构建阶段（builder）中复制文件，/app/main 是构建阶段中编译好的二进制文件路径。
COPY --from=builder ${WORKDIR}/main .

# 运行应用程序
CMD ["./main"]
