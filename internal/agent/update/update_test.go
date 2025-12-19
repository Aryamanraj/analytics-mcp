package update

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestVerifyManifest(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	raw := []byte(`{"version":"1.0.0"}`)
	sig := ed25519.Sign(priv, raw)

	if err := VerifyManifest(raw, sig, base64.StdEncoding.EncodeToString(pub)); err != nil {
		t.Fatalf("expected verification to succeed, got %v", err)
	}

	if err := VerifyManifest(raw, sig, "invalid"); err == nil {
		t.Fatalf("expected failure for invalid pubkey")
	}

	wrongPub, _, _ := ed25519.GenerateKey(rand.Reader)
	if err := VerifyManifest(raw, sig, base64.StdEncoding.EncodeToString(wrongPub)); err == nil {
		t.Fatalf("expected verification failure with wrong pubkey")
	}
}

func TestFetchManifest(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	manifest := Manifest{
		Name:       "payram-analytics",
		Channel:    "stable",
		Version:    "1.2.3",
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

	gotManifest, gotRaw, gotSig, err := FetchManifest(context.Background(), srv.URL, "stable")
	if err != nil {
		t.Fatalf("FetchManifest error: %v", err)
	}

	if gotManifest.Version != manifest.Version || gotManifest.Name != manifest.Name {
		t.Fatalf("unexpected manifest: %+v", gotManifest)
	}
	if string(gotRaw) != string(raw) {
		t.Fatalf("raw mismatch")
	}
	if string(gotSig) != string(sig) {
		t.Fatalf("sig mismatch")
	}

	if err := VerifyManifest(gotRaw, gotSig, base64.StdEncoding.EncodeToString(pub)); err != nil {
		t.Fatalf("verification failed: %v", err)
	}
}
