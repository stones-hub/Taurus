package mcps

import (
	"Taurus/pkg/mcp/mcp_server"
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func init() {
	Router(mcp_server.Core.GetHandler())
}

func Router(handler *server.MCPServer) {
	log.Println("mcp router init")

	handler.AddTool(
		mcp.NewTool(
			"Echo",                                // tool name
			mcp.WithDescription("Echo the input"), // tool description
			mcp.WithString("input",
				mcp.Description("The input to echo"),
				mcp.Required(),
			), // input parameter
			mcp.WithString("output",
				mcp.Description("The output to echo"),
				mcp.Required(),
			), // output parameter
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

			arguments := request.Params.Arguments
			input := arguments["input"].(string)
			output := arguments["output"].(string)

			if input == "" || output == "" {
				log.Println("input or output is empty")
			}

			fmt.Printf("input: %s, output: %s\n", input, output)

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("input: %s, output: %s", input, output),
					},
				},
			}, nil
		},
	)

	handler.AddTool(
		mcp.NewTool(
			"乘法",                      // tool name
			mcp.WithDescription("乘法"), // tool description
			mcp.WithNumber("input1",
				mcp.Description("输入一个数字"),
				mcp.Required(),
			), // input parameter
			mcp.WithNumber("input2",
				mcp.Description("输入一个数字"),
				mcp.Required(),
			), // output parameter
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

			arguments := request.Params.Arguments
			input1 := arguments["input1"].(float64)
			input2 := arguments["input2"].(float64)

			if input1 == 0 || input2 == 0 {
				log.Println("input1 or input2 is empty")
			}

			fmt.Printf("input1: %f, input2: %f", input1, input2)

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("input1: %f, input2: %f, result: %f", input1, input2, input1*input2),
					},
				},
			}, nil
		},
	)

	handler.AddTool(
		mcp.NewTool(
			"Add",                                  // tool name
			mcp.WithDescription("Add two numbers"), // tool description
			mcp.WithNumber("input1",
				mcp.Description("The first number"),
				mcp.Required(),
			), // input parameter
			mcp.WithNumber("input2",
				mcp.Description("The second number"),
				mcp.Required(),
			), // output parameter
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

			arguments := request.Params.Arguments
			input1 := arguments["input1"].(float64)
			input2 := arguments["input2"].(float64)

			if input1 == 0 || input2 == 0 {
				log.Println("input1 or input2 is empty")
			}

			fmt.Printf("input1: %f, input2: %f", input1, input2)

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("input1: %f, input2: %f", input1, input2),
					},
				},
			}, nil
		},
	)

}
