package main

import (
	"Taurus/internal"
	"Taurus/internal/app"
	"Taurus/internal/controller"
	"Taurus/internal/customware"
	"Taurus/pkg/middleware"
	"Taurus/pkg/router"
	"net/http"
)

func main() {
	router.AddRouter(router.Router{
		Path: "/health",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}),
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
				Handler: http.HandlerFunc(internal.Core.DemoCtrl.Get),
				Middleware: []router.MiddlewareFunc{
					// middleware.ApiKeyAuthMiddleware,
					customware.HostMiddleware,
					middleware.ValidationMiddleware(&controller.DemoRequest{}),
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
			customware.HostMiddleware,
		},
	})

	app.Default()
}
