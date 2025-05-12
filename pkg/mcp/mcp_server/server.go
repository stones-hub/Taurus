package mcp_server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

const (
	TransportStdio = "stdio"
	TransportSSE   = "sse"
)

type ServerOption func(*Server)

type Server struct {
	Name             string
	Version          string
	Addr             string // for sse
	Transport        string // set to "stdio" or "sse"
	Handler          *server.MCPServer
	StdioErrorLogger *log.Logger             // for stdio
	StdioContextFunc server.StdioContextFunc // for stdio
	SSEContextFunc   server.SSEContextFunc   // for sse

	stdio *server.StdioServer
	sse   *server.SSEServer
}

// initialize server config
type ServerConfig struct {
	Name        string
	Version     string
	Addr        string
	Transport   string
	Subscribe   bool // 订阅
	ListChanged bool // 列表变化
	Prompt      bool // 提示
	Tool        bool // 工具
}

var (
	Core *Server
)

func InitializeServer(config *ServerConfig, opts ...ServerOption) *Server {

	if config.Addr == "" {
		config.Addr = "localhost:8080"
	}

	if config.Transport == "" {
		config.Transport = TransportStdio
	}

	if config.Name == "" {
		config.Name = "taurus-mcp-server"
	}

	if config.Version == "" {
		config.Version = "0.0.1"
	}

	s := &Server{
		Name:      config.Name,
		Version:   config.Version,
		Addr:      config.Addr,
		Transport: config.Transport,
		Handler: server.NewMCPServer(config.Name, config.Version,
			server.WithResourceCapabilities(config.Subscribe, config.ListChanged),
			server.WithPromptCapabilities(config.Prompt),
			server.WithToolCapabilities(config.Tool),
			server.WithLogging(),
		),
	}
	for _, opt := range opts {
		opt(s)
	}

	Core = s

	return s
}

func WithErrorLogger(logger *log.Logger) ServerOption {
	return func(s *Server) {
		s.StdioErrorLogger = logger
	}
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

func WithTransport(transport string) ServerOption {
	return func(s *Server) {
		s.Transport = transport
	}
}

func WithHandler(handler *server.MCPServer) ServerOption {
	return func(s *Server) {
		s.Handler = handler
	}
}

func WithAddr(addr string) ServerOption {
	return func(s *Server) {
		s.Addr = addr
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

func (s *Server) ListenAndServe() error {

	if s.Addr == "" || s.Transport == "" || s.Handler == nil {
		return fmt.Errorf("invalid server config")
	}
	log.Printf("model protocol context server : %s, %s\n", s.Name, s.Version)

	if s.Transport == TransportStdio {
		s.stdio = server.NewStdioServer(s.Handler)
		if s.StdioErrorLogger == nil {
			s.StdioErrorLogger = defaultLogger()
		}
		s.stdio.SetErrorLogger(s.StdioErrorLogger)
		if s.StdioContextFunc == nil {
			s.StdioContextFunc = defaultStdioContextFunc
		}
		s.stdio.SetContextFunc(s.StdioContextFunc)
		go func() {
			// 启动stdio，不应该设置何时取消context, 因为我们会在程序中用信号集中管理优雅的退出
			err := s.stdio.Listen(context.Background(), os.Stdin, os.Stdout)
			if err != nil && err != io.EOF && err != context.Canceled {
				log.Fatalf("stdio server listen error : %s\n", err.Error())
			}
			log.Println("stdio server goroutine over.")
		}()

	} else if s.Transport == TransportSSE {
		if s.SSEContextFunc == nil {
			s.SSEContextFunc = defaultSSEContextFunc
		}
		s.sse = server.NewSSEServer(s.Handler,
			server.WithBaseURL(fmt.Sprintf("http://%s", s.Addr)),
			server.WithSSEContextFunc(s.SSEContextFunc),
		)
		port := strings.Split(s.Addr, ":")[1]
		go func() { // 启动sse
			err := s.sse.Start(fmt.Sprintf(":%s", port))
			if err != nil {
				log.Fatalf("sse server listen error : %s\n", err.Error())
			}
			log.Println("sse server goroutine over.")
		}()
	} else {
		return fmt.Errorf("invalid transport: %s", s.Transport)
	}

	return nil
}

func (s *Server) Shutdown() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if s.sse != nil {
			s.sse.Shutdown(ctx)
		}
	}()
}

func defaultLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags)
}

func defaultSSEContextFunc(ctx context.Context, r *http.Request) context.Context {
	fmt.Println("sse context")
	return ctx
}

func defaultStdioContextFunc(ctx context.Context) context.Context {
	fmt.Println("stdio context")
	return ctx
}

/*
{
  "mcpServers": {
    "mcp-server-stdio": {
      "disabled": false,
      "timeout": 30,
      "command": "/Users/yelei/data/code/projects/go/Taurus/release/taurus-v0.0.1/taurus",
      "args" : [
        "-config",
        "/Users/yelei/data/code/projects/go/Taurus/release/taurus-v0.0.1/config"
      ],
      "transportType": "stdio"
    },
    "mcp-server-sse": {
      "autoApprove": [
        "Add",
        "Echo",
        "Echo2"
      ],
      "disabled": true,
      "timeout": 60,
      "url": "http://127.0.0.1:9001/sse",
      "transportType": "sse"
    }
  }
}
*/
