package mcp_server

import (
	"fmt"
	"log"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

const (
	TransportStdio = "stdio"
	TransportSSE   = "sse"
)

type ServerOption func(*Server)

type Server struct {
	Addr         string
	Transport    string // stdio, sse
	Handler      *server.MCPServer
	SSEContext   server.SSEContextFunc   // sse transport context
	StdioContext server.StdioContextFunc // stdio transport context
}

func WithStdioContextFunc(fn server.StdioContextFunc) ServerOption {
	return func(s *Server) {
		s.StdioContext = fn
	}
}

func WithSSEContextFunc(fn server.SSEContextFunc) ServerOption {
	return func(s *Server) {
		s.SSEContext = fn
	}
}

func WithTransport(transport string) ServerOption {
	return func(s *Server) {
		s.Transport = transport
	}
}

func WithAddr(addr string) ServerOption {
	return func(s *Server) {
		s.Addr = addr
	}
}

func WithHandler(name string, version string) ServerOption {

	handler := server.NewMCPServer(name, version,
		server.WithResourceCapabilities(true, true), // 资源能力
		server.WithPromptCapabilities(true),         // 提示能力
		server.WithToolCapabilities(true),           // 工具能力
		server.WithLogging(),                        // 添加日志)
	)

	return func(s *Server) {
		s.Handler = handler
	}
}

func (s *Server) String() string {
	return fmt.Sprintf("Server{Addr: %s, Transport: %s, SSEContext: %v, StdioContext: %v}", s.Addr, s.Transport, s.SSEContext, s.StdioContext)
}

func (s *Server) ListenAndServe(opts ...ServerOption) error {
	for _, opt := range opts {
		opt(s)
	}
	log.Printf("Server: %s", s)

	if s.Transport == TransportStdio {
		return server.ServeStdio(s.Handler,
			server.WithStdioContextFunc(s.StdioContext))
	} else if s.Transport == TransportSSE {
		sse := server.NewSSEServer(s.Handler,
			server.WithBaseURL(fmt.Sprintf("http://%s", s.Addr)),
			server.WithSSEContextFunc(s.SSEContext),
		)
		port := strings.Split(s.Addr, ":")[1]
		return sse.Start(fmt.Sprintf(":%s", port))
	} else {
		return fmt.Errorf("invalid transport: %s", s.Transport)
	}
}
