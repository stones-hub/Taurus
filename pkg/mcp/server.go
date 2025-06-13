// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

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
	TransportSSE            = "sse"             // 单机和集群都适合，但是在集群下需要维护有状态的session， Nginx同一个请求来源要路由到同一个服务上才可以
	TransportStreamableHTTP = "streamable_http" // 适合集群部署场景
	ModeStateful            = "stateful"        // 保存上下文
	ModeStateless           = "stateless"       // 不保存上下文
)

type MCPServer struct {
	server *server.Server
}

func NewMCPServer(name, version, transportName string, mode string) (*MCPServer, func(), error) {

	var stateMode transport.StateMode
	switch mode {
	case ModeStateful:
		stateMode = transport.Stateful
	case ModeStateless:
		stateMode = transport.Stateless
	default:
		stateMode = transport.Stateful
	}

	mcpTransport, mcpHandler := getTransport(transportName, stateMode)

	mcpServer, err := server.NewServer(mcpTransport, server.WithServerInfo(protocol.Implementation{
		Name:    name,
		Version: version,
	}))

	if err != nil {
		log.Fatal(err)
	}

	switch h := mcpHandler.(type) {
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
		log.Fatal(fmt.Errorf("unknown handler type: %T", mcpHandler))
	}

	GlobalMCPServer = &MCPServer{
		server: mcpServer,
	}
	return GlobalMCPServer, func() {
		if err := GlobalMCPServer.Shutdown(context.Background()); err != nil {
			log.Fatal(fmt.Errorf("failed to shutdown mcp server: %v", err))
		}
	}, nil
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

func (s *MCPServer) registerTool(tool *protocol.Tool, handler server.ToolHandlerFunc) {
	s.server.RegisterTool(tool, handler)
}

func (s *MCPServer) unregisterTool(name string) {
	s.server.UnregisterTool(name)
}

func (s *MCPServer) registerPrompt(prompt *protocol.Prompt, handler server.PromptHandlerFunc) {
	s.server.RegisterPrompt(prompt, handler)
}

func (s *MCPServer) unregisterPrompt(name string) {
	s.server.UnregisterPrompt(name)
}

func (s *MCPServer) registerResource(resource *protocol.Resource, handler server.ResourceHandlerFunc) {
	s.server.RegisterResource(resource, handler)
}

func (s *MCPServer) unregisterResource(name string) {
	s.server.UnregisterResource(name)
}

func (s *MCPServer) registerResourceTemplate(resourceTemplate *protocol.ResourceTemplate, handler server.ResourceHandlerFunc) {
	s.server.RegisterResourceTemplate(resourceTemplate, handler)
}

func (s *MCPServer) unregisterResourceTemplate(name string) {
	s.server.UnregisterResourceTemplate(name)
}

func (s *MCPServer) RegisterHandler(handler *Handler) {
	for _, tool := range handler.GetTools() {
		// s.server.RegisterTool(tool.ToolName, tool.ToolHandler)
		s.registerTool(tool.ToolName, tool.ToolHandler)
	}
	for _, prompt := range handler.GetPrompts() {
		// s.server.RegisterPrompt(prompt.PromptName, prompt.PromptHandler)
		s.registerPrompt(prompt.PromptName, prompt.PromptHandler)
	}
	for _, resource := range handler.GetResources() {
		// s.server.RegisterResource(resource.ResourceName, resource.ResourceHandler)
		s.registerResource(resource.ResourceName, resource.ResourceHandler)
	}
	for _, resourceTemplate := range handler.GetResourceTemplates() {
		// s.server.RegisterResourceTemplate(resourceTemplate.ResourceTemplateName, resourceTemplate.ResourceTemplateHandler)
		s.registerResourceTemplate(resourceTemplate.ResourceTemplateName, resourceTemplate.ResourceTemplateHandler)
	}
}

// for stdio transport, run the server in the main thread
func (s *MCPServer) Run() error {
	return s.server.Run()
}
