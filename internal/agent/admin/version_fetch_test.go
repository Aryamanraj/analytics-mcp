package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/payram/payram-analytics-mcp-server/internal/version"
)

func TestFetchChildVersionSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(version.Info{Version: "v1", Commit: "abc", BuildDate: "today"})
	}))
	defer srv.Close()

	client := &http.Client{Timeout: time.Second}
	res := fetchChildVersion(context.Background(), client, srv.URL)
	if res.Error != nil {
		t.Fatalf("expected no error, got %+v", res.Error)
	}
	if res.Info == nil || res.Info.Version != "v1" || res.Info.Commit != "abc" || res.Info.BuildDate != "today" {
		t.Fatalf("unexpected info: %#v", res.Info)
	}
}

func TestFetchChildVersionFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := &http.Client{Timeout: time.Second}
	res := fetchChildVersion(context.Background(), client, srv.URL)
	if res.Error == nil {
		t.Fatalf("expected error for non-200 response")
	}
	if res.Info != nil {
		t.Fatalf("expected no info on error")
	}
}

func TestFetchChildVersionInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	client := &http.Client{Timeout: time.Second}
	res := fetchChildVersion(context.Background(), client, srv.URL)
	if res.Error == nil {
		t.Fatalf("expected error for invalid JSON")
	}
	if res.Info != nil {
		t.Fatalf("expected no info on error")
	}
}
