package main

import (
	"Taurus/internal/app"
	"Taurus/internal/controller"
	"Taurus/internal/mid"

	"Taurus/pkg/middleware"
	"Taurus/pkg/router"
	"Taurus/pkg/websocket"
	"net/http"

	_ "Taurus/internal/controller/crons" // 没有依赖的包， 包体内的init是不会被执行的的; 所以导入
	_ "Taurus/internal/controller/mcps"  // 没有依赖的包， 包体内的init是不会被执行的的; 所以导入
	_ "Taurus/internal/log_formatter"    // 没有依赖的包， 包体内的init是不会被执行的的; 所以导入
)

func main() {
	router.AddRouter(router.Router{
		Path: "/health",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}),
	})

	router.AddRouter(router.Router{
		Path: "/ws",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			websocket.HandleWebSocket(w, r, app.GlobalInjector.DemoWs.HandleMessage)
		}),
		Middleware: []router.MiddlewareFunc{
			middleware.ErrorHandlerMiddleware,
			middleware.TraceMiddleware,
		},
	})

	router.AddRouterGroup(router.RouteGroup{
		Prefix: "/v1/api",
		Middleware: []router.MiddlewareFunc{
			middleware.CorsMiddleware,
			middleware.ErrorHandlerMiddleware,
			middleware.TraceMiddleware,
		},
		Routes: []router.Router{
			{
				Path:    "/",
				Handler: http.HandlerFunc(app.GlobalInjector.DemoCtrl.Get),
				Middleware: []router.MiddlewareFunc{
					// middleware.ApiKeyAuthMiddleware,
					mid.HostMiddleware,
				},
			},
		},
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

	router.AddRouter(router.Router{
		Path:    "/",
		Handler: http.HandlerFunc(controller.ServeMarkdownDoc),
		Middleware: []router.MiddlewareFunc{
			mid.HostMiddleware,
		},
	})

	app.Default()
}
