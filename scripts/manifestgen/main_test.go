package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/payram/payram-analytics-mcp-server/internal/agent/update"
)

func TestGenerateWritesManifestAndSignature(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	dir := t.TempDir()
	opts := Options{
		Name:       "payram-analytics",
		Channel:    "stable",
		Version:    "v1.2.3",
		Notes:      "test release",
		ReleasedAt: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
		ChatURL:    "https://example.com/chat.bin",
		ChatSHA:    "abc123",
		MCPURL:     "https://example.com/mcp.bin",
		MCPSHA:     "def456",
		CoreMin:    "1.0.0",
		CoreMax:    "1.2.x",
		OutputDir:  dir,
		PrivKeyB64: base64.StdEncoding.EncodeToString(priv),
	}

	raw, sig, pubB64, err := Generate(opts)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	manifestPath := filepath.Join(dir, "manifest.json")
	sigPath := manifestPath + ".sig"

	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("manifest not written: %v", err)
	}
	if _, err := os.Stat(sigPath); err != nil {
		t.Fatalf("sig not written: %v", err)
	}

	var manifest update.Manifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if manifest.Version != "1.2.3" {
		t.Fatalf("version mismatch: %s", manifest.Version)
	}
	if manifest.Artifacts.Chat.URL == "" || manifest.Artifacts.MCP.URL == "" {
		t.Fatalf("artifact urls missing")
	}

	if err := update.VerifyManifest(raw, sig, pubB64); err != nil {
		t.Fatalf("verify: %v", err)
	}
}
