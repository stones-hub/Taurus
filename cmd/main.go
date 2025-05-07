package main

import (
	"Taurus/internal/app"

	"Taurus/pkg/middleware"
	"Taurus/pkg/router"
	"Taurus/pkg/websocket"
	"net/http"

	_ "Taurus/internal/controller/crons" // 没有依赖的包， 包体内的init是不会被执行的的; 所以导入
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
					middleware.ApiKeyAuthMiddleware,
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

	router.AddRouter(router.Router{
		Path: "/",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/static/index.html", http.StatusFound)
		}),
		Middleware: []router.MiddlewareFunc{},
	})

	app.Default()
}

/*
- http.Dir(path):
	职责: 创建一个文件系统，决定了文件系统的位置（path的位置），规定path的地方才可以被访问, 其他地方不可以被访问

- http.FileServer(http.Dir(path)):
	职责: 创建一个文件服务器，解析http请求，并从文件系统中(http.Dir(path))读取文件，返回给客户端

- http.StripPrefix(prefix, handler):
	职责: 从请求路径中移除指定的前缀(prefix)后，将请求转发给下一个处理器(handler)

注意：其实http.FileServer(http.Dir(path)) 就可以直接解决静态文件的问题，但是有弊端，比如:
1. http.Handle("/", http.FileServer(http.Dir("/project/static"))) ,
	当访问 /css/style.css 时，FileServer 会查找 /project/static/css/style.css 但问题是，这会把整个网站根目录映射到静态文件目录，所有URL都会尝试查找静态文件
2. http.Handle("/static/", http.FileServer(http.Dir("/project/static")))
	当访问 /static/css/style.css 时，FileServer 会查找 /project/static/static/css/style.css 多了一个static
所以让三者配合起来一起使用,让您可以更灵活地组织网站结构，而不受文件系统结构的限制

- http.Redirect(w, r, url, code):
	职责: 重定向请求, 支持的重定向状态有以下几种
	1. http.StatusMovedPermanently (301) 永久重定向，告诉浏览器和搜索引擎该资源已永久移动
	2. http.StatusFound (302) 临时重定向，告诉浏览器和搜索引擎该资源已临时移动
	3. http.StatusSeeOther (303) "查看其他"，通常用于POST请求后重定向到GET页面
	4. http.StatusTemporaryRedirect (307) 临时重定向，保留原请求方法
	5. http.StatusPermanentRedirect (308) 永久重定向，保留原请求方法

	注意：
	1. 重定向状态码(code) 是可选的，默认是 http.StatusFound (302)
*/
