package supervisor

import (
	"testing"

	"github.com/payram/payram-analytics-mcp-server/internal/agent/secrets"
)

func TestChildEnvInjectsOpenAIKey(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)
	if err := secrets.PutOpenAIKey(home, "sk-secret"); err != nil {
		t.Fatalf("put secret: %v", err)
	}

	c := newChild("chat", "echo", nil, Config{BufferLines: 10, InitialBackoff: 1, MaxBackoff: 1, TerminateTimeout: 1})
	env := c.childEnv()
	found := false
	for _, kv := range env {
		if kv == "OPENAI_API_KEY=sk-secret" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected OPENAI_API_KEY injected")
	}
}

func TestChildEnvDoesNotOverrideExistingOpenAIKey(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)
	if err := secrets.PutOpenAIKey(home, "sk-secret"); err != nil {
		t.Fatalf("put secret: %v", err)
	}
	t.Setenv("OPENAI_API_KEY", "from-env")

	c := newChild("chat", "echo", nil, Config{BufferLines: 10, InitialBackoff: 1, MaxBackoff: 1, TerminateTimeout: 1})
	env := c.childEnv()
	foundSecret := false
	foundEnv := false
	for _, kv := range env {
		if kv == "OPENAI_API_KEY=sk-secret" {
			foundSecret = true
		}
		if kv == "OPENAI_API_KEY=from-env" {
			foundEnv = true
		}
	}
	if !foundEnv {
		t.Fatalf("expected existing OPENAI_API_KEY kept")
	}
	if foundSecret {
		t.Fatalf("secret should not override existing env")
	}
}
