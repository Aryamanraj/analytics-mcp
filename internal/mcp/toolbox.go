package mcp

import (
	"context"
	"encoding/json"

	"github.com/payram/payram-analytics-mcp-server/internal/protocol"
)

// Tool defines the behavior of a single MCP tool.
type Tool interface {
	Descriptor() protocol.ToolDescriptor
	Invoke(ctx context.Context, raw json.RawMessage) (protocol.CallResult, *protocol.ResponseError)
}

// Toolbox stores and dispatches tools by name.
type Toolbox struct {
	tools map[string]Tool
}

// NewToolbox constructs a toolbox with the provided tools.
func NewToolbox(tools ...Tool) *Toolbox {
	m := make(map[string]Tool, len(tools))
	for _, t := range tools {
		desc := t.Descriptor()
		m[desc.Name] = t
	}
	return &Toolbox{tools: m}
}

// Describe returns all tool descriptors.
func (tb *Toolbox) Describe() []protocol.ToolDescriptor {
	list := make([]protocol.ToolDescriptor, 0, len(tb.tools))
	for _, t := range tb.tools {
		list = append(list, t.Descriptor())
	}
	return list
}

// Call invokes a named tool.
func (tb *Toolbox) Call(ctx context.Context, name string, args json.RawMessage) (protocol.CallResult, *protocol.ResponseError) {
	tool, ok := tb.tools[name]
	if !ok {
		return protocol.CallResult{}, &protocol.ResponseError{Code: -32601, Message: "tool not found"}
	}
	return tool.Invoke(ctx, args)
}
