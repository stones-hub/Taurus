.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/calculator.proto

.PHONY: run-server
run-server:
	go run server/main.go

.PHONY: run-client
run-client:
	go run client/main.go 