package chatapi

// OpenAI-compatible request/response shapes (subset).

type ChatCompletionRequest struct {
	Model       string          `json:"model"`
	Messages    []OAChatMessage `json:"messages"`
	Tools       []OATool        `json:"tools,omitempty"`
	ToolChoice  interface{}     `json:"tool_choice,omitempty"`
	Temperature *float64        `json:"temperature,omitempty"`
}

type OAChatMessage struct {
	Role       string       `json:"role"`
	Content    string       `json:"content,omitempty"`
	Name       string       `json:"name,omitempty"`
	ToolCallID string       `json:"tool_call_id,omitempty"`
	ToolCalls  []OAToolCall `json:"tool_calls,omitempty"`
}

type OATool struct {
	Type     string     `json:"type"`
	Function OAFunction `json:"function"`
}

type OAFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

type OAToolCall struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Function OAToolCallFunc `json:"function"`
}

type OAToolCallFunc struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Model   string                 `json:"model"`
	Choices []ChatChoice           `json:"choices"`
	Usage   map[string]interface{} `json:"usage,omitempty"`
}

type ChatChoice struct {
	Index        int           `json:"index"`
	Message      OAChatMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}
