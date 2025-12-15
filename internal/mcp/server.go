package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/payram/payram-analytics-mcp-server/internal/protocol"
)

// Server handles MCP JSON-RPC requests against a toolbox.
type Server struct {
	toolbox *Toolbox
}

// NewServer wires a toolbox into an MCP server.
func NewServer(tb *Toolbox) *Server {
	return &Server{toolbox: tb}
}

// Handle routes a single request.
func (s *Server) Handle(ctx context.Context, req protocol.Request) (protocol.Response, error) {
	if err := validateJSONRPC(req); err != nil {
		return protocol.Response{JSONRPC: "2.0", ID: normalizeID(req.ID), Error: err}, nil
	}

	switch req.Method {
	case "initialize":
		return protocol.Response{JSONRPC: "2.0", ID: normalizeID(req.ID), Result: map[string]any{
			"protocolVersion": "2024-11-05",
			"serverInfo": map[string]string{
				"name":    "payram-analytics-mcp-server",
				"version": "0.1.0",
			},
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
		}}, nil
	case "ping":
		return protocol.Response{JSONRPC: "2.0", ID: normalizeID(req.ID), Result: map[string]any{}}, nil
	case "tools/list":
		return protocol.Response{JSONRPC: "2.0", ID: normalizeID(req.ID), Result: protocol.ListResult{Tools: s.toolbox.Describe()}}, nil
	case "tools/call":
		var params protocol.CallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return protocol.Response{JSONRPC: "2.0", ID: normalizeID(req.ID), Error: &protocol.ResponseError{Code: -32602, Message: "invalid params"}}, nil
		}
		if params.Name == "" {
			return protocol.Response{JSONRPC: "2.0", ID: normalizeID(req.ID), Error: &protocol.ResponseError{Code: -32602, Message: "tool name required"}}, nil
		}
		result, toolErr := s.toolbox.Call(ctx, params.Name, params.Args)
		if toolErr != nil {
			return protocol.Response{JSONRPC: "2.0", ID: normalizeID(req.ID), Error: toolErr}, nil
		}
		return protocol.Response{JSONRPC: "2.0", ID: normalizeID(req.ID), Result: result}, nil
	default:
		return protocol.Response{JSONRPC: "2.0", ID: normalizeID(req.ID), Error: &protocol.ResponseError{Code: -32601, Message: "method not found"}}, nil
	}
}

// WriteError builds a response with an error and wraps encode issues.
func WriteError(id any, code int, message string, err error) protocol.Response {
	detail := message
	if err != nil {
		detail = fmt.Sprintf("%s: %v", message, err)
	}
	return protocol.Response{JSONRPC: "2.0", ID: normalizeID(id), Error: &protocol.ResponseError{Code: code, Message: detail}}
}

func validateJSONRPC(req protocol.Request) *protocol.ResponseError {
	if req.JSONRPC != "" && req.JSONRPC != "2.0" {
		return &protocol.ResponseError{Code: -32600, Message: "invalid jsonrpc version"}
	}
	return nil
}

func normalizeID(id any) any {
	if id == nil {
		return "0"
	}
	switch v := id.(type) {
	case string:
		return v
	case float64:
		return v
	case int, int32, int64, uint32, uint64:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
