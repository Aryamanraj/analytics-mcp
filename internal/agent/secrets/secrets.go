package secrets

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/payram/payram-analytics-mcp-server/internal/agent/update"
)

// Secrets holds persisted secret material.
type Secrets struct {
	OpenAIAPIKey string `json:"openai_api_key,omitempty"`
}

// Load returns secrets and their source: "env", "state", or "missing".
func Load(home string) (Secrets, string, error) {
	if home == "" {
		home = update.HomeDir()
	}

	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return Secrets{OpenAIAPIKey: key}, "env", nil
	}

	path := pathFor(home)
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Secrets{}, "missing", nil
		}
		return Secrets{}, "", err
	}

	var s Secrets
	if err := json.Unmarshal(raw, &s); err != nil {
		return Secrets{}, "", err
	}

	if s.OpenAIAPIKey != "" {
		return s, "state", nil
	}
	return Secrets{}, "missing", nil
}

// PutOpenAIKey writes the key atomically with 0600 permissions.
func PutOpenAIKey(home, key string) error {
	if key == "" {
		return fmt.Errorf("openai api key empty")
	}
	if home == "" {
		home = update.HomeDir()
	}

	dir := filepath.Join(home, "state")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	path := pathFor(home)
	tmp := path + ".tmp"

	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}

	enc, err := json.Marshal(Secrets{OpenAIAPIKey: key})
	if err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return err
	}

	if _, err := f.Write(enc); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	// best-effort fsync on directory
	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		_ = d.Close()
	}

	return nil
}

// DeleteOpenAIKey removes the stored key.
func DeleteOpenAIKey(home string) error {
	if home == "" {
		home = update.HomeDir()
	}
	path := pathFor(home)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func pathFor(home string) string {
	return filepath.Join(home, "state", "secrets.json")
}
