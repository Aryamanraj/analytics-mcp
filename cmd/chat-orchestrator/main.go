package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/payram/payram-analytics-mcp-server/internal/chatserver"
)

func main() {
	port := envOr("CHAT_PORT", "3000")
	mcpURL := envOr("MCP_SERVER_URL", "http://localhost:8080/")
	staticDir := envOr("CHAT_STATIC_DIR", "web")
	apiKey := envOr("OPENAI_API_KEY", "")
	model := envOr("OPENAI_MODEL", "gpt-4o-mini")
	baseURL := envOr("OPENAI_BASE_URL", "https://api.openai.com/v1")

	flag.StringVar(&port, "port", port, "port to listen on")
	flag.StringVar(&mcpURL, "mcp", mcpURL, "base URL for MCP server (HTTP)")
	flag.StringVar(&staticDir, "static", staticDir, "directory for static assets")
	flag.StringVar(&apiKey, "openai-key", apiKey, "OpenAI API key")
	flag.StringVar(&model, "openai-model", model, "OpenAI model")
	flag.StringVar(&baseURL, "openai-base", baseURL, "OpenAI base URL")
	flag.Parse()

	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is required for the chat orchestrator")
	}

	mcpClient := chatserver.NewMCPClient(mcpURL)
	llmClient := chatserver.NewLLMClient(apiKey, model, baseURL)
	srv := chatserver.NewChatServer(mcpClient, llmClient, staticDir)

	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	handler := logging(mux)
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Chat orchestrator listening on %s (MCP: %s, Model: %s)", addr, mcpURL, model)

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}
