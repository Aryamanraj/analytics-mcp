package admin

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/payram/payram-analytics-mcp-server/internal/agent/supervisor"
	"github.com/payram/payram-analytics-mcp-server/internal/agent/update"
)

func TestUpdateAvailableSuccess(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	manifest := update.Manifest{
		Name:       "payram-analytics",
		Channel:    "stable",
		Version:    "1.7.4",
		Notes:      "bug fixes",
		ReleasedAt: time.Date(2025, 12, 18, 0, 0, 0, 0, time.UTC),
	}

	raw, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	sig := ed25519.Sign(priv, raw)

	mux := http.NewServeMux()
	mux.HandleFunc("/stable/manifest.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(raw)
	})
	mux.HandleFunc("/stable/manifest.json.sig", func(w http.ResponseWriter, _ *http.Request) {
		w.Write(sig)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	t.Setenv("PAYRAM_AGENT_ADMIN_TOKEN", "tok")
	t.Setenv("PAYRAM_AGENT_ADMIN_ALLOWLIST", "")
	t.Setenv("PAYRAM_AGENT_UPDATE_BASE_URL", srv.URL)
	t.Setenv("PAYRAM_AGENT_UPDATE_PUBKEY_B64", base64.StdEncoding.EncodeToString(pub))

	sup := &supervisor.Supervisor{}
	handler := NewMux(sup)

	req := httptest.NewRequest(http.MethodGet, "/admin/update/available?channel=stable", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Set("Authorization", "Bearer tok")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("missing data: %v", body)
	}

	if data["target_version"] != manifest.Version {
		t.Fatalf("unexpected target_version: %v", data["target_version"])
	}
	if data["notes"] != manifest.Notes {
		t.Fatalf("unexpected notes: %v", data["notes"])
	}
}

func TestUpdateAvailableSignatureInvalid(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	manifest := update.Manifest{Version: "1.0.0"}
	raw, _ := json.Marshal(manifest)
	sig := ed25519.Sign(priv, raw)

	mux := http.NewServeMux()
	mux.HandleFunc("/stable/manifest.json", func(w http.ResponseWriter, _ *http.Request) { w.Write(raw) })
	mux.HandleFunc("/stable/manifest.json.sig", func(w http.ResponseWriter, _ *http.Request) { w.Write(sig) })

	srv := httptest.NewServer(mux)
	defer srv.Close()

	wrongPub, _, _ := ed25519.GenerateKey(rand.Reader)

	t.Setenv("PAYRAM_AGENT_ADMIN_TOKEN", "tok")
	t.Setenv("PAYRAM_AGENT_ADMIN_ALLOWLIST", "")
	t.Setenv("PAYRAM_AGENT_UPDATE_BASE_URL", srv.URL)
	t.Setenv("PAYRAM_AGENT_UPDATE_PUBKEY_B64", base64.StdEncoding.EncodeToString(wrongPub))

	sup := &supervisor.Supervisor{}
	handler := NewMux(sup)

	req := httptest.NewRequest(http.MethodGet, "/admin/update/available", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Set("Authorization", "Bearer tok")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	errObj, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object: %v", body)
	}
	if errObj["code"] != "SIGNATURE_INVALID" {
		t.Fatalf("unexpected error code: %v", errObj["code"])
	}
}
