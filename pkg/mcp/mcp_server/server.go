package mcp_server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	TransportStdio = "stdio"
	TransportSSE   = "sse"
)

var (
	Core *Server
)

type ServerOption func(*Server)

type Server struct {
	Name             string
	Version          string
	Addr             string                  // for sse
	Transport        string                  // set to "stdio" or "sse"
	StdioErrorLogger *log.Logger             // for stdio
	StdioContextFunc server.StdioContextFunc // for stdio
	SSEContextFunc   server.SSEContextFunc   // for sse
	handler          *server.MCPServer
}

func (s *Server) GetHandler() *server.MCPServer {
	if s.handler == nil {
		log.Fatal("MCP server handler is nil")
		return nil
	}
	return s.handler
}

func setupHooks(hooks *server.Hooks) {
	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		fmt.Printf("beforeAny: %s, %v, %v\n", method, id, message)
	})
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		fmt.Printf("onSuccess: %s, %v, %v, %v\n", method, id, message, result)
	})
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		fmt.Printf("onError: %s, %v, %v, %v\n", method, id, message, err)
	})
	hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
		fmt.Printf("beforeInitialize: %v, %v\n", id, message)
	})
	hooks.AddOnRequestInitialization(func(ctx context.Context, id any, message any) error {
		fmt.Printf("AddOnRequestInitialization: %v, %v\n", id, message)
		// authorization verification and other preprocessing tasks are performed.
		return nil
	})
	hooks.AddAfterInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest, result *mcp.InitializeResult) {
		fmt.Printf("afterInitialize: %v, %v, %v\n", id, message, result)
	})
	hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		fmt.Printf("afterCallTool: %v, %v, %v\n", id, message, result)
	})
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		fmt.Printf("beforeCallTool: %v, %v\n", id, message)
	})
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}

	hooks := &server.Hooks{}
	setupHooks(hooks)

	s.handler = server.NewMCPServer(
		s.Name,
		s.Version,
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(hooks),
	)
	return s
}

func WithName(name string) ServerOption {
	return func(s *Server) {
		s.Name = name
	}
}

func WithVersion(version string) ServerOption {
	return func(s *Server) {
		s.Version = version
	}
}

func WithAddr(addr string) ServerOption {
	return func(s *Server) {
		s.Addr = addr
	}
}

func WithTransport(transport string) ServerOption {
	return func(s *Server) {
		s.Transport = transport
	}
}

func WithStdioErrorLogger(logger *log.Logger) ServerOption {
	return func(s *Server) {
		s.StdioErrorLogger = logger
	}
}

func WithStdioContextFunc(fn server.StdioContextFunc) ServerOption {
	return func(s *Server) {
		s.StdioContextFunc = fn
	}
}

func WithSSEContextFunc(fn server.SSEContextFunc) ServerOption {
	return func(s *Server) {
		s.SSEContextFunc = fn
	}
}

func (s *Server) String() string {
	return fmt.Sprintf("Server{Addr: %s, Transport: %s}", s.Addr, s.Transport)
}

func defaultLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags)
}

func defaultSSEContextFunc(ctx context.Context, r *http.Request) context.Context {
	// TODO: 添加SSE上下文, 可用于验证请求来源
	// 获取请求头中的Authorization
	log.Printf("SSE request header: %v", r.Header)
	return ctx
}

func defaultStdioContextFunc(ctx context.Context) context.Context {
	// TODO: 添加stdio上下文, 可用于验证请求来源
	log.Printf("stdio context: %v", ctx)
	return ctx
}

func (s *Server) ListenAndServe() {
	if s.Transport == TransportSSE {

		if s.SSEContextFunc == nil {
			s.SSEContextFunc = defaultSSEContextFunc
		}

		sseServer := server.NewSSEServer(s.handler,
			server.WithBaseURL("http://"+s.Addr),
			server.WithSSEContextFunc(s.SSEContextFunc),
		)

		log.Printf("SSE server listening on %s", s.Addr)
		// localhost:8080 , 获取 :8080部分，去掉localhost
		port := strings.Split(s.Addr, ":")[1]
		if err := sseServer.Start(":" + port); err != nil {
			log.Fatalf("MCP SSE Server error: %v", err)
		}
	} else if s.Transport == TransportStdio {
		if s.StdioContextFunc == nil {
			s.StdioContextFunc = defaultStdioContextFunc
		}

		if s.StdioErrorLogger == nil {
			s.StdioErrorLogger = defaultLogger()
		}

		if err := server.ServeStdio(s.handler,
			server.WithStdioContextFunc(s.StdioContextFunc),
			server.WithErrorLogger(s.StdioErrorLogger),
		); err != nil {
			log.Fatalf("MCP stdio Server error: %v", err)
		}
	} else {
		log.Fatalf("MCP Server Invalid transport: %s", s.Transport)
	}
}

/*
{
  "mcpServers": {
    "mcp-server-stdio": {
      "autoApprove": [
        "Add",
        "Echo",
        "乘法"
      ],
      "disabled": false,
      "timeout": 30,
      "command": "/Users/yelei/data/code/projects/go/Taurus/release/taurus-v0.0.1/taurus",
      "args": [
        "-config",
        "/Users/yelei/data/code/projects/go/Taurus/release/taurus-v0.0.1/config",
        "-env",
        "/Users/yelei/data/code/projects/go/Taurus/release/taurus-v0.0.1/.env.local"
      ],
      "transportType": "stdio"
    },
    "mcp-server-sse": {
      "autoApprove": [
        "Echo",
        "Echo2",
        "Add",
        "乘法"
      ],
      "disabled": false,
      "timeout": 30,
      "url": "http://127.0.0.1:9001/sse",
      "transportType": "sse"
    }
  }
}
*/
