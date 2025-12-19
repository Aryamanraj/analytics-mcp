package update

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func writeTempBinary(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write temp binary: %v", err)
	}
	return path
}

func TestEnsureSeedReleaseCreatesBaseline(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)

	srcDir := t.TempDir()
	chatSrc := writeTempBinary(t, srcDir, "chat-src", "chatbin")
	mcpSrc := writeTempBinary(t, srcDir, "mcp-src", "mcpbin")
	t.Setenv("PAYRAM_AGENT_SEED_CHAT_SRC", chatSrc)
	t.Setenv("PAYRAM_AGENT_SEED_MCP_SRC", mcpSrc)

	seeded, version, err := EnsureSeedRelease(context.Background(), home)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	if !seeded {
		t.Fatalf("expected seeded true")
	}
	if version != "0.0.0" {
		t.Fatalf("version mismatch: %s", version)
	}

	releaseDir := filepath.Join(home, "releases", "0.0.0")
	chatBin := filepath.Join(releaseDir, chatBinaryName)
	mcpBin := filepath.Join(releaseDir, mcpBinaryName)
	if _, err := os.Stat(chatBin); err != nil {
		t.Fatalf("chat not copied: %v", err)
	}
	if _, err := os.Stat(mcpBin); err != nil {
		t.Fatalf("mcp not copied: %v", err)
	}

	if target, err := os.Readlink(filepath.Join(releaseDir, "chat")); err != nil || filepath.Base(target) != chatBinaryName {
		t.Fatalf("chat compat link invalid: %v target=%s", err, target)
	}
	if target, err := os.Readlink(filepath.Join(releaseDir, "mcp")); err != nil || filepath.Base(target) != mcpBinaryName {
		t.Fatalf("mcp compat link invalid: %v target=%s", err, target)
	}

	curTarget, err := os.Readlink(CurrentSymlink())
	if err != nil {
		t.Fatalf("current symlink: %v", err)
	}
	if curTarget != filepath.Join(home, "releases", "0.0.0") {
		t.Fatalf("current target mismatch: %s", curTarget)
	}
}

func TestEnsureSeedReleaseIdempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)

	srcDir := t.TempDir()
	chatSrc := writeTempBinary(t, srcDir, "chat-src", "chatbin")
	mcpSrc := writeTempBinary(t, srcDir, "mcp-src", "mcpbin")
	t.Setenv("PAYRAM_AGENT_SEED_CHAT_SRC", chatSrc)
	t.Setenv("PAYRAM_AGENT_SEED_MCP_SRC", mcpSrc)

	if seeded, _, err := EnsureSeedRelease(context.Background(), home); err != nil || !seeded {
		t.Fatalf("initial seed failed: seeded=%v err=%v", seeded, err)
	}

	seeded, version, err := EnsureSeedRelease(context.Background(), home)
	if err != nil {
		t.Fatalf("second seed: %v", err)
	}
	if seeded {
		t.Fatalf("expected seeded false on second call")
	}
	if version != "" {
		t.Fatalf("expected empty version on second call, got %s", version)
	}
}

func TestEnsureSeedReleaseMissingSources(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)
	t.Setenv("PAYRAM_AGENT_SEED_CHAT_SRC", filepath.Join(home, "nochat"))
	t.Setenv("PAYRAM_AGENT_SEED_MCP_SRC", filepath.Join(home, "nomcp"))

	if _, _, err := EnsureSeedRelease(context.Background(), home); err == nil {
		t.Fatalf("expected error when sources missing")
	}
}
