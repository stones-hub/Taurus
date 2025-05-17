package resources

import (
	"Taurus/pkg/logx"
	"Taurus/pkg/mcp"
	"context"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func init() {
	// add resource
	mcp.MCPHandler.RegisterResource(&protocol.Resource{
		URI:      "file:///test.txt",
		Name:     "test.txt",
		MimeType: "text/plain-txt",
	}, TestResource)
}

func TestResource(ctx context.Context, request *protocol.ReadResourceRequest) (*protocol.ReadResourceResult, error) {
	logx.Core.Info("default", "read resource, request: %v", request)
	return &protocol.ReadResourceResult{
		Contents: []protocol.ResourceContents{
			protocol.TextResourceContents{
				URI:      "file:///test.txt",
				MimeType: "text/plain-txt",
				Text:     "test",
			},
		},
	}, nil
}
