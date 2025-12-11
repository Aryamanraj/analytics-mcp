package chatserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/payram/payram-analytics-mcp-server/internal/protocol"
)

// ChatServer provides a simple chat-like API that can call MCP tools.
type ChatServer struct {
	mcp       *MCPClient
	llm       *LLMClient
	staticDir string
}

// NewChatServer wires MCP client and static assets location.
func NewChatServer(mcp *MCPClient, llm *LLMClient, staticDir string) *ChatServer {
	return &ChatServer{mcp: mcp, llm: llm, staticDir: staticDir}
}

// RegisterRoutes attaches handlers to the mux.
func (s *ChatServer) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/chat", s.handleChat)
	mux.HandleFunc("/api/tools", s.handleTools)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	fs := http.FileServer(http.Dir(s.staticDir))
	mux.Handle("/", fs)
}

// ChatRequest is the payload from the UI.
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse is what the UI renders.
type ChatResponse struct {
	Reply      string            `json:"reply"`
	ToolCall   *ToolCall         `json:"toolCall,omitempty"`
	ToolResult *ToolResult       `json:"toolResult,omitempty"`
	Error      string            `json:"error,omitempty"`
	Meta       map[string]string `json:"meta,omitempty"`
}

// ToolCall captures an invoked tool.
type ToolCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args,omitempty"`
}

// ToolResult includes a human-readable string.
type ToolResult struct {
	Text string `json:"text"`
}

func (s *ChatServer) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		writeJSON(w, ChatResponse{Error: "message is required"}, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tools, err := s.mcp.ListTools(ctx)
	if err != nil {
		writeJSON(w, ChatResponse{Error: fmt.Sprintf("tools error: %v", err)}, http.StatusBadGateway)
		return
	}

	decision, err := s.llm.Decide(ctx, msg, tools)
	if err != nil {
		writeJSON(w, ChatResponse{Error: fmt.Sprintf("llm error: %v", err)}, http.StatusBadGateway)
		return
	}

	var (
		reply      string
		toolCall   *ToolCall
		toolResult *ToolResult
	)

	switch decision.Action {
	case "tool_call":
		reply = "Using a tool to gather context..."
		toolCall = &ToolCall{Name: decision.Name, Args: decision.Args}
		result, err := s.mcp.CallTool(ctx, decision.Name, decision.Args)
		if err != nil {
			writeJSON(w, ChatResponse{Error: fmt.Sprintf("tool error: %v", err)}, http.StatusBadGateway)
			return
		}
		text := renderContent(result)
		toolResult = &ToolResult{Text: text}
		if decision.Message != "" {
			reply = decision.Message
		} else {
			reply = fmt.Sprintf("Used tool '%s'. Here's what I found:\n%s", decision.Name, text)
		}
	case "respond":
		reply = decision.Message
	default:
		writeJSON(w, ChatResponse{Error: "invalid llm action"}, http.StatusBadGateway)
		return
	}

	writeJSON(w, ChatResponse{
		Reply:      reply,
		ToolCall:   toolCall,
		ToolResult: toolResult,
		Meta: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	}, http.StatusOK)
}

func (s *ChatServer) handleTools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tools, err := s.mcp.ListTools(r.Context())
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadGateway)
		return
	}

	writeJSON(w, tools, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json error: %v", err)
	}
}

// renderContent flattens tool output into readable text.
func renderContent(result protocol.CallResult) string {
	var sb strings.Builder
	for i, c := range result.Content {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(c.Text)
	}
	return strings.TrimSpace(sb.String())
}
