package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/payram/payram-analytics-mcp-server/internal/agent/supervisor"
	"github.com/payram/payram-analytics-mcp-server/internal/version"
)

func NewMux(sup *supervisor.Supervisor) http.Handler {
	if sup == nil {
		panic("supervisor must not be nil")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/version", versionHandler)

	adminGuard := NewAdminMiddlewareFromEnv()
	mux.Handle("/admin/version", adminGuard(http.HandlerFunc(adminVersionHandler)))
	mux.Handle("/admin/child/restart", adminGuard(http.HandlerFunc(restartHandler(sup))))
	mux.Handle("/admin/child/status", adminGuard(http.HandlerFunc(statusHandler(sup))))
	mux.Handle("/admin/logs", adminGuard(http.HandlerFunc(logsHandler(sup))))

	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	RespondOK(w, http.StatusOK, map[string]any{"status": "healthy"})
}

func versionHandler(w http.ResponseWriter, _ *http.Request) {
	RespondOK(w, http.StatusOK, version.Get())
}

func adminVersionHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{Timeout: 2 * time.Second}
	ctx := r.Context()

	chatPort := envPort("PAYRAM_CHAT_PORT", 2358)
	mcpPort := envPort("PAYRAM_MCP_PORT", 3333)

	chatURL := fmt.Sprintf("http://127.0.0.1:%d/version", chatPort)
	mcpURL := fmt.Sprintf("http://127.0.0.1:%d/version", mcpPort)

	chat := fetchChildVersion(ctx, client, chatURL)
	mcp := fetchChildVersion(ctx, client, mcpURL)

	RespondOK(w, http.StatusOK, map[string]any{
		"agent": version.Get(),
		"chat":  chat,
		"mcp":   mcp,
	})
}

func restartHandler(sup *supervisor.Supervisor) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RespondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only POST allowed")
			return
		}

		if err := sup.RestartAll(); err != nil {
			RespondError(w, http.StatusInternalServerError, "RESTART_FAILED", err.Error())
			return
		}

		RespondOK(w, http.StatusOK, map[string]any{"status": "restarted"})
	}
}

func statusHandler(sup *supervisor.Supervisor) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		RespondOK(w, http.StatusOK, sup.Status())
	}
}

func logsHandler(sup *supervisor.Supervisor) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		component := r.URL.Query().Get("component")
		if component == "" {
			RespondError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "component is required")
			return
		}

		tail := 200
		if rawTail := r.URL.Query().Get("tail"); rawTail != "" {
			if parsed, err := strconv.Atoi(rawTail); err == nil && parsed > 0 {
				tail = parsed
			}
		}

		logs := sup.Logs(component, tail)
		if logs == nil {
			RespondError(w, http.StatusBadRequest, "INVALID_COMPONENT", "component must be chat or mcp")
			return
		}

		RespondOK(w, http.StatusOK, map[string]any{
			"component": component,
			"lines":     logs,
		})
	}
}

type childVersionResult struct {
	Info  *version.Info `json:"info,omitempty"`
	Error *respError    `json:"error,omitempty"`
}

func fetchChildVersion(ctx context.Context, client *http.Client, url string) childVersionResult {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return childVersionResult{Error: &respError{Code: "FETCH_FAILED", Message: err.Error()}}
	}

	resp, err := client.Do(req)
	if err != nil {
		return childVersionResult{Error: &respError{Code: "FETCH_FAILED", Message: err.Error()}}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return childVersionResult{Error: &respError{Code: "FETCH_FAILED", Message: fmt.Sprintf("status %d", resp.StatusCode)}}
	}

	var info version.Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return childVersionResult{Error: &respError{Code: "FETCH_FAILED", Message: err.Error()}}
	}

	return childVersionResult{Info: &info}
}

func envPort(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			return p
		}
	}
	return fallback
}
