package tools

import (
	"Taurus/pkg/mcp"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func init() {
	mcp.MCPHandler.RegisterTool(CurrentTimeTool(), CurrentTime)
}

func CurrentTime(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(CurrentTimeReq)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("parse timezone with error: %v", err)
	}
	text := fmt.Sprintf(`current time is %s`, time.Now().In(loc))

	log.Printf("texti11111: %s", text)

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}, nil
}

func CurrentTimeTool() *protocol.Tool {
	tool, err := protocol.NewTool("current_time", "Get current time with timezone, Asia/Shanghai is default", CurrentTimeReq{})
	if err != nil {
		panic(fmt.Sprintf("Failed to create tool: %v", err))
	}
	return tool
}

type CurrentTimeReq struct {
	Timezone string `json:"timezone" description:"current time timezone"`
}
