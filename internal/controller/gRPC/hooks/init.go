package hooks

/*
 * ------------------------------------------------------------
 * 用于初始化gRPC的中间件、拦截器、验证器
 * ------------------------------------------------------------
 */

import (
	"Taurus/pkg/grpc/server"
	"Taurus/pkg/grpc/server/interceptor"
	"Taurus/pkg/grpc/server/middleware"
	"Taurus/pkg/telemetry"
	"log"
	"sync"
)

var once sync.Once

func InitgRPCHooks() {

	once.Do(func() {
		log.Printf("-------------------------------- 初始化 gRPC 中间件 --------------------------------")
		tracer := telemetry.Provider.Tracer("gRPC-server")
		server.RegisterMiddleware(HostMiddleware())
		// System default middleware
		// MetricsMiddleware
		server.RegisterMiddleware(middleware.MetricsMiddleware(tracer))
		// LoggingMiddleware
		server.RegisterMiddleware(middleware.LoggingMiddleware())

		log.Printf("-------------------------------- 初始化 gRPC 拦截器 --------------------------------")
		server.RegisterInterceptor(AuthInterceptor())

		// System default interceptor
		// LoggingServerInterceptor
		server.RegisterInterceptor(interceptor.LoggingServerInterceptor())
		// RecoveryServerInterceptor
		server.RegisterInterceptor(interceptor.RecoveryServerInterceptor())
		// AuthServerInterceptor
		server.RegisterInterceptor(interceptor.AuthServerInterceptor("Bearer 123456"))
		server.RegisterStreamInterceptor(interceptor.AuthStreamServerInterceptor("Bearer 123456"))
		// RateLimitServerInterceptor
		// server.RegisterInterceptor(interceptor.RateLimitServerInterceptor(10))

		// 验证器
		// ValidatorServerInterceptor
		server.RegisterInterceptor(interceptor.UnaryServerValidationInterceptor())
		// ValidatorStreamServerInterceptor
		server.RegisterStreamInterceptor(interceptor.StreamServerValidationInterceptor())

		server.RegisterStreamInterceptor(interceptor.MetricsStreamInterceptor(tracer))
	})
}
