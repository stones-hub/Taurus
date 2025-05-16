package resources

import (
	"Taurus/pkg/logx"
	"Taurus/pkg/mcp"
	"context"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func init() {
	Do()
}

func Do() {
	// add resource
	testResource := &protocol.Resource{
		URI:      "file:///test.txt",
		Name:     "test.txt",
		MimeType: "text/plain-txt",
	}
	testResourceContent := protocol.TextResourceContents{
		URI:      testResource.URI,
		MimeType: testResource.MimeType,
		Text:     "test",
	}

	mcp.MCPHandler.RegisterResource(testResource, func(ctx context.Context, request *protocol.ReadResourceRequest) (*protocol.ReadResourceResult, error) {
		logx.Core.Info("default", "read resource, request: %v", request)
		return &protocol.ReadResourceResult{
			Contents: []protocol.ResourceContents{
				testResourceContent,
			},
		}, nil
	})

}
