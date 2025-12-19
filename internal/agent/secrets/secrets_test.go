package secrets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissing(t *testing.T) {
	home := t.TempDir()
	s, source, err := Load(home)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if source != "missing" {
		t.Fatalf("expected missing source, got %s", source)
	}
	if s.OpenAIAPIKey != "" {
		t.Fatalf("expected empty key")
	}
}

func TestPutLoadDelete(t *testing.T) {
	home := t.TempDir()

	if err := PutOpenAIKey(home, "sk-test"); err != nil {
		t.Fatalf("put: %v", err)
	}

	path := filepath.Join(home, "state", "secrets.json")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected 0600 perms, got %v", info.Mode().Perm())
	}

	s, source, err := Load(home)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if source != "state" {
		t.Fatalf("expected state source, got %s", source)
	}
	if s.OpenAIAPIKey != "sk-test" {
		t.Fatalf("key mismatch: %s", s.OpenAIAPIKey)
	}

	if err := DeleteOpenAIKey(home); err != nil {
		t.Fatalf("delete: %v", err)
	}
	s, source, err = Load(home)
	if err != nil {
		t.Fatalf("load after delete: %v", err)
	}
	if source != "missing" || s.OpenAIAPIKey != "" {
		t.Fatalf("expected missing after delete")
	}
}

func TestEnvBeatsState(t *testing.T) {
	home := t.TempDir()
	if err := PutOpenAIKey(home, "sk-state"); err != nil {
		t.Fatalf("put: %v", err)
	}
	t.Setenv("OPENAI_API_KEY", "sk-env")
	s, source, err := Load(home)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if source != "env" {
		t.Fatalf("expected env source, got %s", source)
	}
	if s.OpenAIAPIKey != "sk-env" {
		t.Fatalf("env key not returned")
	}
}
