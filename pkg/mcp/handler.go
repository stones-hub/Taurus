package mcp

import (
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

var (
	MCPHandler = &Handler{}
)

type Tool struct {
	ToolName    *protocol.Tool
	ToolHandler server.ToolHandlerFunc
}

type Prompt struct {
	PromptName    *protocol.Prompt
	PromptHandler server.PromptHandlerFunc
}

type Resource struct {
	ResourceName    *protocol.Resource
	ResourceHandler server.ResourceHandlerFunc
}

type ResourceTemplate struct {
	ResourceTemplateName    *protocol.ResourceTemplate
	ResourceTemplateHandler server.ResourceHandlerFunc
}

type Handler struct {
	tools             []Tool
	prompts           []Prompt
	resources         []Resource
	resourceTemplates []ResourceTemplate
}

func (h *Handler) RegisterTool(tool *protocol.Tool, handler server.ToolHandlerFunc) {
	h.tools = append(h.tools, Tool{
		ToolName:    tool,
		ToolHandler: handler,
	})
}

func (h *Handler) RegisterPrompt(prompt *protocol.Prompt, handler server.PromptHandlerFunc) {
	h.prompts = append(h.prompts, Prompt{
		PromptName:    prompt,
		PromptHandler: handler,
	})
}

func (h *Handler) RegisterResource(resource *protocol.Resource, handler server.ResourceHandlerFunc) {
	h.resources = append(h.resources, Resource{
		ResourceName:    resource,
		ResourceHandler: handler,
	})
}

func (h *Handler) RegisterResourceTemplate(resourceTemplate *protocol.ResourceTemplate, handler server.ResourceHandlerFunc) {
	h.resourceTemplates = append(h.resourceTemplates, ResourceTemplate{
		ResourceTemplateName:    resourceTemplate,
		ResourceTemplateHandler: handler,
	})
}

func (h *Handler) UnregisterTool(name string) {
	for i, tool := range h.tools {
		if tool.ToolName.Name == name {
			h.tools = append(h.tools[:i], h.tools[i+1:]...)
		}
	}
	if GlobalMCPServer != nil {
		GlobalMCPServer.unregisterTool(name)
	}
}

func (h *Handler) UnregisterPrompt(name string) {
	for i, prompt := range h.prompts {
		if prompt.PromptName.Name == name {
			h.prompts = append(h.prompts[:i], h.prompts[i+1:]...)
		}
	}
	if GlobalMCPServer != nil {
		GlobalMCPServer.unregisterPrompt(name)
	}
}

func (h *Handler) UnregisterResource(name string) {
	for i, resource := range h.resources {
		if resource.ResourceName.Name == name {
			h.resources = append(h.resources[:i], h.resources[i+1:]...)
		}
	}
	if GlobalMCPServer != nil {
		GlobalMCPServer.unregisterResource(name)
	}
}

func (h *Handler) UnregisterResourceTemplate(name string) {
	for i, resourceTemplate := range h.resourceTemplates {
		if resourceTemplate.ResourceTemplateName.Name == name {
			h.resourceTemplates = append(h.resourceTemplates[:i], h.resourceTemplates[i+1:]...)
		}
	}
	if GlobalMCPServer != nil {
		GlobalMCPServer.unregisterResourceTemplate(name)
	}
}

func (h *Handler) GetTools() []Tool {
	return h.tools
}

func (h *Handler) GetPrompts() []Prompt {
	return h.prompts
}

func (h *Handler) GetResources() []Resource {
	return h.resources
}

func (h *Handler) GetResourceTemplates() []ResourceTemplate {
	return h.resourceTemplates
}
