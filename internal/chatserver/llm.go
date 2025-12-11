package chatserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/payram/payram-analytics-mcp-server/internal/protocol"
)

// LLMClient calls a chat-completions style API (OpenAI-compatible) to decide how to respond.
type LLMClient struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// LLMDecision is the structured action returned from the LLM.
type LLMDecision struct {
	Action  string         `json:"action"`            // "respond" or "tool_call"
	Message string         `json:"message,omitempty"` // when action == respond
	Name    string         `json:"name,omitempty"`    // when action == tool_call
	Args    map[string]any `json:"args,omitempty"`    // when action == tool_call
}

// NewLLMClient configures an OpenAI-compatible caller.
func NewLLMClient(apiKey, model, baseURL string) *LLMClient {
	url := strings.TrimRight(baseURL, "/")
	if url == "" {
		url = "https://api.openai.com/v1"
	}
	return &LLMClient{
		apiKey:  apiKey,
		model:   model,
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Decide asks the LLM whether to call a tool or respond directly.
func (c *LLMClient) Decide(ctx context.Context, userMessage string, tools []protocol.ToolDescriptor) (LLMDecision, error) {
	var decision LLMDecision

	if c.apiKey == "" {
		return decision, errors.New("missing LLM API key")
	}

	prompt := buildSystemPrompt(tools)
	reqBody := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: prompt},
			{Role: "user", Content: userMessage},
		},
		Temperature:    0.2,
		ResponseFormat: map[string]string{"type": "json_object"},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return decision, fmt.Errorf("encode llm request: %w", err)
	}

	endpoint := c.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return decision, fmt.Errorf("build llm request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return decision, fmt.Errorf("call llm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return decision, fmt.Errorf("llm returned status %d", resp.StatusCode)
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return decision, fmt.Errorf("decode llm response: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return decision, errors.New("llm returned no choices")
	}

	content := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	parsed, err := parseDecision(content)
	if err != nil {
		return decision, fmt.Errorf("parse llm decision: %w", err)
	}

	return parsed, nil
}

func buildSystemPrompt(tools []protocol.ToolDescriptor) string {
	var b strings.Builder
	b.WriteString("You are PayRam's chat orchestrator. Use available tools when they match the user's request.\n")
	b.WriteString("Available tools (name: description):\n")
	for _, t := range tools {
		b.WriteString("- ")
		b.WriteString(t.Name)
		if t.Description != "" {
			b.WriteString(": ")
			b.WriteString(t.Description)
		}
		b.WriteString("\n")
	}
	b.WriteString("\nOutput ONLY compact JSON. Formats:\n")
	b.WriteString("{\"action\":\"tool_call\",\"name\":\"tool_name\",\"args\":{}}\n")
	b.WriteString("or\n")
	b.WriteString("{\"action\":\"respond\",\"message\":\"your reply\"}\n")
	b.WriteString("Use tool_call whenever a tool directly helps answer. Keep args object even if empty.")
	return b.String()
}

func parseDecision(raw string) (LLMDecision, error) {
	var dec LLMDecision
	raw = strings.TrimSpace(stripCodeFence(raw))
	if raw == "" {
		return dec, errors.New("empty response")
	}
	if err := json.Unmarshal([]byte(raw), &dec); err != nil {
		return dec, err
	}
	switch dec.Action {
	case "respond":
		if dec.Message == "" {
			return dec, errors.New("respond action missing message")
		}
	case "tool_call":
		if dec.Name == "" {
			return dec, errors.New("tool_call missing name")
		}
		if dec.Args == nil {
			dec.Args = map[string]any{}
		}
	default:
		return dec, fmt.Errorf("unknown action %q", dec.Action)
	}
	return dec, nil
}

func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimPrefix(s, "json")
		s = strings.TrimSpace(s)
		if idx := strings.LastIndex(s, "```"); idx >= 0 {
			s = s[:idx]
		}
	}
	return s
}

// Minimal OpenAI-style request/response payloads

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model          string            `json:"model"`
	Messages       []chatMessage     `json:"messages"`
	Temperature    float64           `json:"temperature"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}
