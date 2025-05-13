package mcp

import (
	"Taurus/pkg/router"
	"fmt"
	"log"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

func NewMCPServer(name, version, transportName string, stateMode transport.StateMode) *server.Server {

	t, handler := getTransport(transportName, stateMode)

	mcpServer, err := server.NewServer(t, server.WithServerInfo(protocol.Implementation{
		Name:    name,
		Version: version,
	}))

	if err != nil {
		log.Fatal(err)
	}

	switch h := handler.(type) {
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
		log.Fatal(fmt.Errorf("unknown handler type: %T", handler))
	}

	return mcpServer

}

func getTransport(transportName string, stateMode transport.StateMode) (transport.ServerTransport, interface{}) {
	var err error
	var t transport.ServerTransport
	var handler interface{}

	switch transportName {
	case "stdio":
		log.Println("start current time mcp server with stdio transport")
		t = transport.NewStdioServerTransport()
	case "sse":
		log.Println("start current time mcp server with sse transport")
		var sseHandler *transport.SSEHandler
		t, sseHandler, err = transport.NewSSEServerTransportAndHandler("/message")
		if err != nil {
			log.Fatal(fmt.Errorf("failed to create sse transport: %v", err))
		}
		fmt.Println("sseHandler:---------> ", sseHandler)
		handler = sseHandler
	case "streamable_http":
		log.Println("start current time mcp server with streamable http transport")
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
