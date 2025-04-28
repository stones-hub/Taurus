package rpc

import (
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"time"
)

func NewGrpcClient(host string, port int, timeout int) (err error, conn *grpc.ClientConn) {
	addr := fmt.Sprintf("%s:%d", host, port)
	creds := grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12}))
	keepAlive := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Duration(timeout) * time.Second,
		Timeout:             time.Duration(timeout) * time.Second,
		PermitWithoutStream: true,
	})
	conn, err = grpc.NewClient(addr, creds, keepAlive)
	return
}
