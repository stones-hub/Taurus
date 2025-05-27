package main

import (
	"Taurus/internal"
	"Taurus/internal/app"
	"Taurus/internal/controller"
	"Taurus/internal/customware"
	"Taurus/pkg/middleware"
	"Taurus/pkg/router"
	"Taurus/pkg/telemetry"
	"net/http"
)

func main() {

	// 测试validate
	router.AddRouterGroup(router.RouteGroup{
		Prefix: "/v1/api",
		Middleware: []router.MiddlewareFunc{
			middleware.CorsMiddleware,
			middleware.ErrorHandlerMiddleware,
		},
		Routes: []router.Router{
			{
				Path:    "/",
				Handler: http.HandlerFunc(internal.Core.ValidateCtrl.TestValidateMiddleware),
				Middleware: []router.MiddlewareFunc{
					customware.HostMiddleware,
					middleware.ValidationMiddleware(&controller.ValidateRequest{}), // 验证请求是否符合ValidateRequest结构体
				},
			},
		},
	})

	tracer := telemetry.Provider.Tracer("http-server")

	// 测试trace中间件
	router.AddRouter(router.Router{
		Path:    "/trace",
		Handler: http.HandlerFunc(internal.Core.TraceCtrl.TestTraceMiddleware),
		Middleware: []router.MiddlewareFunc{
			middleware.TraceMiddleware(tracer),
			customware.HostMiddleware,
		},
	})

	// 设置健康检查
	router.AddRouter(router.Router{
		Path: "/health",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}),
	})

	// 设置静态文件地址
	router.AddRouter(router.Router{
		Path: "/static/",
		// 浏览器访问 http://localhost:8080/static/css/style.css 首先会去掉/static/ 剩下css/style.css, 然后去相对于应用根目录 ./static去找css/style.css文件返回
		Handler:    http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
		Middleware: []router.MiddlewareFunc{},
	})

	// 重定向到静态文件
	router.AddRouter(router.Router{
		Path: "/demo",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/static/index.html", http.StatusFound)
		}),
		Middleware: []router.MiddlewareFunc{},
	})

	// 设置markdown文档
	router.AddRouter(router.Router{
		Path:    "/",
		Handler: http.HandlerFunc(controller.ServeMarkdownDoc),
		Middleware: []router.MiddlewareFunc{
			customware.HostMiddleware,
		},
	})

	app.Default()
}
