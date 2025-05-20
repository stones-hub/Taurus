package server

import (
	"google.golang.org/grpc"
)

// ServiceRegistrar 服务注册接口
type ServiceRegistrar interface {
	RegisterService(server *grpc.Server)
}

// 服务注册表
var serviceRegistry = make(map[string]ServiceRegistrar)

// RegisterService 注册服务
func RegisterService(name string, service ServiceRegistrar) {
	serviceRegistry[name] = service
}

// GetRegisteredServices 获取所有注册的服务
func GetRegisteredServices() map[string]ServiceRegistrar {
	return serviceRegistry
}
