package server

import (
	"Taurus/pkg/grpc/attributes"

	"google.golang.org/grpc"
)

// ServiceRegistrar 服务注册接口
type ServiceRegistrar interface {
	RegisterService(server *grpc.Server)
}

// 服务注册表
var (
	serviceRegistry          = make(map[string]ServiceRegistrar)
	serviceMiddleware        = make([]attributes.UnaryMiddleware, 0)
	serviceStreamMiddleware  = make([]attributes.StreamMiddleware, 0)
	serviceInterceptor       = make([]grpc.UnaryServerInterceptor, 0)
	serviceStreamInterceptor = make([]grpc.StreamServerInterceptor, 0)
)

// RegisterService 注册服务
func RegisterService(name string, service ServiceRegistrar) {
	serviceRegistry[name] = service
}

// GetRegisteredServices 获取所有注册的服务
func GetRegisteredServices() map[string]ServiceRegistrar {
	return serviceRegistry
}

func RegisterMiddleware(middleware attributes.UnaryMiddleware) {
	serviceMiddleware = append(serviceMiddleware, middleware)
}

func RegisterStreamMiddleware(middleware attributes.StreamMiddleware) {
	serviceStreamMiddleware = append(serviceStreamMiddleware, middleware)
}

func RegisterInterceptor(interceptor grpc.UnaryServerInterceptor) {
	serviceInterceptor = append(serviceInterceptor, interceptor)
}

func RegisterStreamInterceptor(interceptor grpc.StreamServerInterceptor) {
	serviceStreamInterceptor = append(serviceStreamInterceptor, interceptor)
}

func GetServiceMiddleware() []attributes.UnaryMiddleware {
	return serviceMiddleware
}

func GetServiceStreamMiddleware() []attributes.StreamMiddleware {
	return serviceStreamMiddleware
}

func GetServiceInterceptor() []grpc.UnaryServerInterceptor {
	return serviceInterceptor
}

func GetServiceStreamInterceptor() []grpc.StreamServerInterceptor {
	return serviceStreamInterceptor
}
