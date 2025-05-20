package server

import (
	"log"

	"google.golang.org/grpc"
)

type ServiceRegister interface {
	Register(server *grpc.Server)
}

var (
	// 服务注册表
	serviceRegistry = make(map[string]ServiceRegister)
)

// 将服务写入注册表
func RegisterService(name string, service ServiceRegister) {
	if _, ok := serviceRegistry[name]; ok {
		log.Printf("service %s already registered", name)
	}
	serviceRegistry[name] = service
}

// 获取注册表
func GetServiceRegistry() map[string]ServiceRegister {
	return serviceRegistry
}
