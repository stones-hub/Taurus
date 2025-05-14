package mcp

import (
	"Taurus/pkg/router"
	"context"
	"fmt"
	"log"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

var (
	GlobalMCPServer *MCPServer
)

const (
	TransportStdio          = "stdio"           // 适合单机部署场景
	TransportSSE            = "sse"             // 适合单机部署场景, 需要维护有状态的session
	TransportStreamableHTTP = "streamable_http" // 适合集群部署场景
	ModeStateful            = "stateful"        // 保存上下文
	ModeStateless           = "stateless"       // 不保存上下文
)

type MCPServer struct {
	server *server.Server
}

func NewMCPServer(name, version, transportName string, mode string) *MCPServer {

	var stateMode transport.StateMode
	switch mode {
	case ModeStateful:
		stateMode = transport.Stateful
	case ModeStateless:
		stateMode = transport.Stateless
	default:
		stateMode = transport.Stateful
	}

	t, handler := getTransport(transportName, stateMode)

	mcpServer, err := server.NewServer(t, server.WithServerInfo(protocol.Implementation{
		Name:    name,
		Version: version,
	}))

	if err != nil {
		log.Fatal(err)
	}

	switch h := handler.(type) {
	case *transport.SSEHandler:
		router.AddRouter(router.Router{
			Path:       "/sse",
			Handler:    h.HandleSSE(),
			Middleware: nil,
		})

		router.AddRouter(router.Router{
			Path:       "/message",
			Handler:    h.HandleMessage(),
			Middleware: nil,
		})
	case *transport.StreamableHTTPHandler:
		router.AddRouter(router.Router{
			Path:       "/mcp",
			Handler:    h.HandleMCP(),
			Middleware: nil,
		})
	default:
		log.Fatal(fmt.Errorf("unknown handler type: %T", handler))
	}

	GlobalMCPServer = &MCPServer{
		server: mcpServer,
	}
	return GlobalMCPServer
}

func (s *MCPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func getTransport(transportName string, stateMode transport.StateMode) (transport.ServerTransport, interface{}) {
	var err error
	var t transport.ServerTransport
	var handler interface{}

	switch transportName {
	case TransportStdio:
		log.Println("start mcp server with stdio transport")
		t = transport.NewStdioServerTransport()
	case TransportSSE:
		log.Println("start mcp server with sse transport")
		var sseHandler *transport.SSEHandler
		t, sseHandler, err = transport.NewSSEServerTransportAndHandler("/message")
		if err != nil {
			log.Fatal(fmt.Errorf("failed to create sse transport: %v", err))
		}
		handler = sseHandler
	case TransportStreamableHTTP:
		log.Println("start mcp server with streamable http transport")
		var streamableHandler *transport.StreamableHTTPHandler
		t, streamableHandler, err = transport.NewStreamableHTTPServerTransportAndHandler(transport.WithStreamableHTTPServerTransportAndHandlerOptionStateMode(stateMode))
		if err != nil {
			log.Fatal(fmt.Errorf("failed to create streamable http transport: %v", err))
		}
		handler = streamableHandler
	default:
		log.Fatal(fmt.Errorf("unknown transport name: %s", transportName))
	}

	return t, handler
}

func (s *MCPServer) RegisterTool(tool *protocol.Tool, handler server.ToolHandlerFunc) {
	s.server.RegisterTool(tool, handler)
}

func (s *MCPServer) UnregisterTool(name string) {
	s.server.UnregisterTool(name)
}

func (s *MCPServer) RegisterPrompt(prompt *protocol.Prompt, handler server.PromptHandlerFunc) {
	s.server.RegisterPrompt(prompt, handler)
}

func (s *MCPServer) UnregisterPrompt(name string) {
	s.server.UnregisterPrompt(name)
}

func (s *MCPServer) RegisterResource(resource *protocol.Resource, handler server.ResourceHandlerFunc) {
	s.server.RegisterResource(resource, handler)
}

func (s *MCPServer) UnregisterResource(name string) {
	s.server.UnregisterResource(name)
}

func (s *MCPServer) RegisterResourceTemplate(resourceTemplate *protocol.ResourceTemplate, handler server.ResourceHandlerFunc) {
	s.server.RegisterResourceTemplate(resourceTemplate, handler)
}

func (s *MCPServer) UnregisterResourceTemplate(name string) {
	s.server.UnregisterResourceTemplate(name)
}

func (s *MCPServer) RegisterHandler(handler *Handler) {
	for _, tool := range handler.GetTools() {
		s.server.RegisterTool(tool.ToolName, tool.ToolHandler)
	}
	for _, prompt := range handler.GetPrompts() {
		s.server.RegisterPrompt(prompt.PromptName, prompt.PromptHandler)
	}
	for _, resource := range handler.GetResources() {
		s.server.RegisterResource(resource.ResourceName, resource.ResourceHandler)
	}
	for _, resourceTemplate := range handler.GetResourceTemplates() {
		s.server.RegisterResourceTemplate(resourceTemplate.ResourceTemplateName, resourceTemplate.ResourceTemplateHandler)
	}
}

// for stdio transport, run the server in the main thread
func (s *MCPServer) Run() error {
	return s.server.Run()
}
