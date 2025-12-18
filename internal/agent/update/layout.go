package update

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const defaultHomeDir = "/var/lib/payram-mcp"

const (
	chatBinaryName = "payram-analytics-chat"
	mcpBinaryName  = "payram-analytics-mcp"
)

// HomeDir resolves the agent home directory from PAYRAM_AGENT_HOME or default.
func HomeDir() string {
	if v := os.Getenv("PAYRAM_AGENT_HOME"); v != "" {
		return v
	}
	return defaultHomeDir
}

// ReleasesDir returns the releases directory.
func ReleasesDir() string { return filepath.Join(HomeDir(), "releases") }

// ReleaseDir returns the directory for a specific version.
func ReleaseDir(version string) string { return filepath.Join(ReleasesDir(), version) }

// StateDir returns the state directory.
func StateDir() string { return filepath.Join(HomeDir(), "state") }

// LockDir returns the lock directory.
func LockDir() string { return filepath.Join(HomeDir(), "lock") }

// LockFilePath returns the update lock path.
func LockFilePath() string { return filepath.Join(LockDir(), "update.lock") }

// CurrentSymlink returns the current symlink path.
func CurrentSymlink() string { return filepath.Join(HomeDir(), "current") }

// PreviousSymlink returns the previous symlink path.
func PreviousSymlink() string { return filepath.Join(HomeDir(), "previous") }

// EnsureBaseDirs ensures required directories exist.
func EnsureBaseDirs() error {
	for _, dir := range []string{ReleasesDir(), StateDir(), LockDir()} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

var ErrUpdateInProgress = errors.New("update already in progress")

// AcquireUpdateLock obtains an exclusive update lock, writing pid/timestamp into the file.
func AcquireUpdateLock() (func() error, error) {
	if err := EnsureBaseDirs(); err != nil {
		return nil, err
	}

	path := LockFilePath()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return nil, ErrUpdateInProgress
		}
		return nil, err
	}

	_, _ = fmt.Fprintf(f, "pid=%d\nstarted=%s\n", os.Getpid(), time.Now().Format(time.RFC3339))
	_ = f.Close()

	unlock := func() error { return os.Remove(path) }
	return unlock, nil
}

// UpdateSymlinks atomically sets current to newTarget and previous to the old current target.
func UpdateSymlinks(newTarget string) (string, error) {
	if err := EnsureBaseDirs(); err != nil {
		return "", err
	}

	current := CurrentSymlink()
	previous := PreviousSymlink()

	oldTarget, _ := os.Readlink(current)

	if oldTarget != "" {
		tmpPrev := previous + ".tmp"
		_ = os.Remove(tmpPrev)
		if err := os.Symlink(oldTarget, tmpPrev); err != nil {
			return oldTarget, err
		}
		if err := os.Rename(tmpPrev, previous); err != nil {
			_ = os.Remove(tmpPrev)
			return oldTarget, err
		}
	}

	tmpCur := current + ".tmp"
	_ = os.Remove(tmpCur)
	if err := os.Symlink(newTarget, tmpCur); err != nil {
		return oldTarget, err
	}
	if err := os.Rename(tmpCur, current); err != nil {
		_ = os.Remove(tmpCur)
		return oldTarget, err
	}

	return oldTarget, nil
}

// VersionFromTarget extracts the version directory name from a symlink target.
func VersionFromTarget(target string) string {
	if target == "" {
		return ""
	}
	return filepath.Base(target)
}

// DefaultChatBin returns the default chat binary path inside the current release.
func DefaultChatBin() string {
	return filepath.Join(CurrentSymlink(), chatBinaryName)
}

// DefaultMCPBin returns the default MCP binary path inside the current release.
func DefaultMCPBin() string {
	return filepath.Join(CurrentSymlink(), mcpBinaryName)
}
