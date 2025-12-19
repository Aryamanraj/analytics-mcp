package update

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateSymlinks(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)

	oldDir := filepath.Join(home, "releases", "1.0.0")
	newDir := filepath.Join(home, "releases", "1.1.0")

	if err := os.MkdirAll(oldDir, 0o755); err != nil {
		t.Fatalf("mkdir old: %v", err)
	}
	if err := os.MkdirAll(newDir, 0o755); err != nil {
		t.Fatalf("mkdir new: %v", err)
	}

	firstOld, err := UpdateSymlinks(oldDir)
	if err != nil {
		t.Fatalf("initial symlink: %v", err)
	}
	if firstOld != "" {
		t.Fatalf("expected empty old target, got %q", firstOld)
	}

	oldTarget, err := UpdateSymlinks(newDir)
	if err != nil {
		t.Fatalf("second symlink: %v", err)
	}
	if oldTarget != oldDir {
		t.Fatalf("expected old target %q got %q", oldDir, oldTarget)
	}

	curTarget, err := os.Readlink(CurrentSymlink())
	if err != nil {
		t.Fatalf("read current: %v", err)
	}
	if curTarget != newDir {
		t.Fatalf("current link mismatch: %q", curTarget)
	}

	prevTarget, err := os.Readlink(PreviousSymlink())
	if err != nil {
		t.Fatalf("read previous: %v", err)
	}
	if prevTarget != oldDir {
		t.Fatalf("previous link mismatch: %q", prevTarget)
	}
}

func TestAcquireUpdateLock(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)

	unlock, err := AcquireUpdateLock()
	if err != nil {
		t.Fatalf("acquire lock: %v", err)
	}

	if _, err := AcquireUpdateLock(); err == nil {
		t.Fatalf("expected in-progress error")
	}

	if err := unlock(); err != nil {
		t.Fatalf("unlock: %v", err)
	}
}
