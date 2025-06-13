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
