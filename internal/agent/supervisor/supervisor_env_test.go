package supervisor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/payram/payram-analytics-mcp-server/internal/agent/secrets"
	"github.com/payram/payram-analytics-mcp-server/internal/agent/update"
)

func TestNewFromEnvUsesCurrentDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)
	t.Setenv("PAYRAM_AGENT_CHAT_BIN", "")
	t.Setenv("PAYRAM_AGENT_MCP_BIN", "")
	current := filepath.Join(home, "current")
	if err := os.MkdirAll(current, 0o755); err != nil {
		t.Fatalf("mkdir current: %v", err)
	}

	chatPath := filepath.Join(current, "payram-analytics-chat")
	mcpPath := filepath.Join(current, "payram-analytics-mcp")
	if err := os.WriteFile(chatPath, []byte("chat"), 0o755); err != nil {
		t.Fatalf("write chat: %v", err)
	}
	if err := os.WriteFile(mcpPath, []byte("mcp"), 0o755); err != nil {
		t.Fatalf("write mcp: %v", err)
	}

	sup, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv error: %v", err)
	}

	if sup.chat.path != update.DefaultChatBin() {
		t.Fatalf("chat path mismatch: %s", sup.chat.path)
	}
	if sup.mcp.path != update.DefaultMCPBin() {
		t.Fatalf("mcp path mismatch: %s", sup.mcp.path)
	}
}

func TestNewFromEnvMissingDefaultFails(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)
	t.Setenv("PAYRAM_AGENT_CHAT_BIN", "")
	t.Setenv("PAYRAM_AGENT_MCP_BIN", "")

	if _, err := NewFromEnv(); err == nil {
		t.Fatalf("expected error for missing binaries")
	}
}

func TestEnsureOpenAIKeyOverridesEmptyEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)
	if err := os.MkdirAll(filepath.Join(home, "state"), 0o755); err != nil {
		t.Fatalf("mkdir state: %v", err)
	}
	if err := secrets.PutOpenAIKey(home, "sk-test-key"); err != nil {
		t.Fatalf("put key: %v", err)
	}

	input := []string{"OPENAI_API_KEY="}
	out := ensureOpenAIKey(input)
	if !hasEnvWithValue(out, "OPENAI_API_KEY") {
		t.Fatalf("expected OPENAI_API_KEY with value from secrets, got %v", out)
	}
}
