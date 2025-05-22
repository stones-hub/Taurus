.PHONY: proto server client

# 生成 proto 文件
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/user.proto

# 运行服务端
server:
	go run server/main.go

# 运行客户端
client:
	go run client/main.go 