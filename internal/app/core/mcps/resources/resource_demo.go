package resources

import (
	"Taurus/pkg/logx"
	"Taurus/pkg/mcp"
	"context"
	"os"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func init() {
	// add resource
	mcp.MCPHandler.RegisterResource(&protocol.Resource{
		URI:      "file:///index.html",
		Name:     "index.html",
		MimeType: "text/plain-txt",
	}, TestResource)
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
