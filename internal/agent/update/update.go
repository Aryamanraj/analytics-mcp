package update

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manifest describes an available update.
type Manifest struct {
	Name          string        `json:"name"`
	Channel       string        `json:"channel"`
	Version       string        `json:"version"`
	ReleasedAt    time.Time     `json:"released_at"`
	Notes         string        `json:"notes"`
	Artifacts     Artifacts     `json:"artifacts"`
	Compatibility Compatibility `json:"compatibility"`
	Revoked       bool          `json:"revoked"`
}

// Artifacts contains binaries for each component.
type Artifacts struct {
	Chat Artifact `json:"chat"`
	MCP  Artifact `json:"mcp"`
}

// Artifact describes a downloadable binary.
type Artifact struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

// Compatibility captures version ranges for dependencies.
type Compatibility struct {
	PayramCore Range `json:"payram_core"`
}

// Range defines min/max versions.
type Range struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

// FetchManifest downloads manifest and signature for a channel.
func FetchManifest(ctx context.Context, baseURL, channel string) (Manifest, []byte, []byte, error) {
	var manifest Manifest
	base := strings.TrimRight(baseURL, "/")
	if channel == "" {
		channel = "stable"
	}
	manifestURL := fmt.Sprintf("%s/%s/manifest.json", base, channel)
	sigURL := manifestURL + ".sig"

	raw, err := fetchBytes(ctx, manifestURL)
	if err != nil {
		return manifest, nil, nil, err
	}

	sig, err := fetchBytes(ctx, sigURL)
	if err != nil {
		return manifest, raw, nil, err
	}

	if err := json.Unmarshal(raw, &manifest); err != nil {
		return manifest, raw, sig, err
	}

	return manifest, raw, sig, nil
}

// VerifyManifest checks a manifest signature using a base64 public key.
func VerifyManifest(raw, sig []byte, pubKeyB64 string) error {
	pub, err := base64.StdEncoding.DecodeString(pubKeyB64)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}
	if len(pub) != ed25519.PublicKeySize {
		return errors.New("invalid public key length")
	}
	if !ed25519.Verify(ed25519.PublicKey(pub), raw, sig) {
		return errors.New("signature verification failed")
	}
	return nil
}

// DownloadToFile streams a URL to a destination file.
func DownloadToFile(ctx context.Context, url, dstPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}

	tmp := dstPath + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(f, resp.Body)
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(tmp)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return closeErr
	}

	if err := os.Rename(tmp, dstPath); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return nil
}

// VerifySHA256 checks a file hash against expected hex.
func VerifySHA256(filePath, expectedHex string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	sum := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(sum, expectedHex) {
		return fmt.Errorf("sha256 mismatch: got %s expected %s", sum, expectedHex)
	}
	return nil
}

func fetchBytes(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
