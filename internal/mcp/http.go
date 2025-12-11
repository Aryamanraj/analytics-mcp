package mcp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/payram/payram-analytics-mcp-server/internal/protocol"
)

// RunHTTP starts an HTTP server that serves MCP JSON-RPC requests via POST.
// Expects a single JSON-RPC request per call. Clients should POST to the root path.
func RunHTTP(server *Server, addr string) error {
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req protocol.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, protocol.Response{Error: &protocol.ResponseError{Code: -32700, Message: "invalid JSON"}}, http.StatusBadRequest)
			return
		}

		resp, err := server.Handle(r.Context(), req)
		if err != nil {
			writeJSON(w, WriteError(req.ID, -32603, "internal error", err), http.StatusInternalServerError)
			return
		}

		writeJSON(w, resp, http.StatusOK)
	})

	log.Printf("HTTP MCP server listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func writeJSON(w http.ResponseWriter, resp protocol.Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(resp)
}
