package protocol

import "encoding/json"

// Request represents a minimal JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string          `json:"jsonrpc,omitempty"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// Response models a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string         `json:"jsonrpc,omitempty"`
	ID      any            `json:"id"`
	Result  any            `json:"result,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

// ResponseError holds JSON-RPC error data.
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ToolDescriptor describes a tool available from the MCP server.
type ToolDescriptor struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema *JSONSchema `json:"inputSchema,omitempty"`
}

// JSONSchema is a minimal subset to describe tool input shapes.
type JSONSchema struct {
	Type                 string                `json:"type,omitempty"`
	Properties           map[string]JSONSchema `json:"properties,omitempty"`
	Required             []string              `json:"required,omitempty"`
	Enum                 []string              `json:"enum,omitempty"`
	Description          string                `json:"description,omitempty"`
	AdditionalProperties any                   `json:"additionalProperties,omitempty"`
}

// ListResult is the payload for tools/list.
type ListResult struct {
	Tools []ToolDescriptor `json:"tools"`
}

// CallParams represents parameters for tools/call.
type CallParams struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"arguments,omitempty"`
}

// ContentPart is a single piece of tool output.
type ContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// CallResult is the payload for a successful tool invocation.
type CallResult struct {
	Content []ContentPart `json:"content"`
}
