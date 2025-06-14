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

package main

import (
	"Taurus/pkg/logx"
	"Taurus/pkg/mcp"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func main() {
	mcp.MCPHandler.RegisterTool(CurrentTimeTool(), CurrentTime)

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

func CurrentTime(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(CurrentTimeReq)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	logx.Core.Info("default", "call tool, request: %v", request)

	loc, err := time.LoadLocation(req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("parse timezone with error: %v", err)
	}
	text := fmt.Sprintf(`current time is %s`, time.Now().In(loc))

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
