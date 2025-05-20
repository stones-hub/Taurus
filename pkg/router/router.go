package router

import (
	"log"
	"net/http"
)

// Router holds the configuration for a route, including its handler and middleware
type Router struct {
	Path       string
	Handler    http.Handler
	Middleware []MiddlewareFunc
}

// RouteGroup holds a group of routes with a common prefix and middleware
type RouteGroup struct {
	Prefix     string
	Middleware []MiddlewareFunc
	Routes     []Router
}

// RouterManager manages all routes and route groups
type RouterManager struct {
	routes          []Router
	routeGroups     []RouteGroup
	registeredPaths map[string]bool // Track registered paths
}

// DefaultManager is the default instance of RouterManager
var DefaultManager = NewRouterManager()

// NewRouterManager creates a new RouterManager
func NewRouterManager() *RouterManager {
	return &RouterManager{
		routes:          []Router{},
		routeGroups:     []RouteGroup{},
		registeredPaths: make(map[string]bool),
	}
}

// AddRouter adds a single route to the manager
func AddRouter(route Router) {
	DefaultManager.routes = append(DefaultManager.routes, route)
}

// AddRouterGroup adds a route group to the manager
func AddRouterGroup(group RouteGroup) {
	DefaultManager.routeGroups = append(DefaultManager.routeGroups, group)
}

// LoadRoutes loads all routes and route groups into a ServeMux
func LoadRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	// Load individual routes
	for _, route := range DefaultManager.routes {
		if DefaultManager.registeredPaths[route.Path] {
			log.Printf("Warning: Path %s is already registered, skipping.\n", route.Path)
			continue
		}
		handler := ChainMiddleware(route.Handler, route.Middleware...)
		mux.Handle(route.Path, handler)
		DefaultManager.registeredPaths[route.Path] = true
	}
	// Load route groups
	for _, group := range DefaultManager.routeGroups {
		for _, route := range group.Routes {
			// Combine group and route middleware, maintaining order
			allMiddleware := append(group.Middleware, route.Middleware...)
			handler := ChainMiddleware(route.Handler, allMiddleware...)
			// Ensure the path is correctly formatted
			fullPath := group.Prefix + route.Path
			if fullPath == "" || DefaultManager.registeredPaths[fullPath] {
				log.Printf("Warning: Path %s is already registered, skipping.\n", fullPath)
				continue // Skip if the full path is empty or already registered
			}
			mux.Handle(fullPath, handler)
			DefaultManager.registeredPaths[fullPath] = true
		}
	}
	return mux
}

// MiddlewareFunc defines a function to process middleware
// It takes an http.Handler and returns an http.Handler
// This allows chaining of middleware functions
type MiddlewareFunc func(http.Handler) http.Handler

// ChainMiddleware applies a list of middleware functions to an http.Handler
func ChainMiddleware(handler http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	// 反转一下，调用顺序不对
	/*
		for _, middleware := range middlewares {
			handler = middleware(handler)
		}
	*/
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
