package main

import (
	"Taurus/pkg/logx"
	"Taurus/pkg/mcp"
	"context"
	"log"
	"os"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func main() {
	// add resource
	mcp.MCPHandler.RegisterResource(&protocol.Resource{
		URI:      "file:///index.html",
		Name:     "index.html",
		MimeType: "text/plain-txt",
	}, TestResource)

	server, _, err := mcp.NewMCPServer("mcp_demo", "1.0.0", "streamable_http", "stateless")
	if err != nil {
		log.Fatalf("Failed to initialize mcp server: %v", err)
	}
	// register handler for mcp server
	server.RegisterHandler(mcp.MCPHandler)

	server.Run()

	defer server.Shutdown(context.Background())

	select {}
}
func TestResource(ctx context.Context, request *protocol.ReadResourceRequest) (*protocol.ReadResourceResult, error) {
	logx.Core.Info("default", "read resource, request: %v", request)

	// 读取html文件
	html, err := os.ReadFile("templates/index.html")
	if err != nil {
		return nil, err
	}

	return &protocol.ReadResourceResult{
		Contents: []protocol.ResourceContents{
			protocol.TextResourceContents{
				URI:      "file:///index.html",
				MimeType: "text/plain-txt",
				Text:     string(html),
			},
		},
	}, nil
}
