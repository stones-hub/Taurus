package prompts

import (
	"Taurus/pkg/logx"
	"Taurus/pkg/mcp"
	"context"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func init() {
	mcp.MCPHandler.RegisterPrompt(prompt(), promptHandler)
}

func prompt() *protocol.Prompt {
	return &protocol.Prompt{
		Name:        "system_prompt",
		Description: "system prompt",
		Arguments: []protocol.PromptArgument{
			{
				Name:        "system_prompt_argument",
				Description: "system prompt argument description",
				Required:    true,
			},
		},
	}
}

func promptHandler(ctx context.Context, request *protocol.GetPromptRequest) (*protocol.GetPromptResult, error) {
	logx.Core.Info("default", "call prompt, request: %v", request)
	return &protocol.GetPromptResult{
		Messages: []protocol.PromptMessage{
			{
				Role: protocol.RoleUser,
				Content: &protocol.TextContent{
					Type: "text",
					Text: "Prompt handler content",
				},
			},
		},
		Description: "Prompt handler description",
	}, nil
}
