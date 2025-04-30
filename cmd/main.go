package main

import (
	"Taurus/internal/app"
	"Taurus/internal/middleware"
	"Taurus/pkg/router"
	"Taurus/pkg/websocket"
	"net/http"
)

func main() {
	router.AddRouter(router.Router{
		Path: "/health",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
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

	app.Default()
}
