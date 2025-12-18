package version

import "testing"

func TestGetDefaults(t *testing.T) {
	origVersion, origCommit, origBuildDate := Version, Commit, BuildDate
	t.Cleanup(func() {
		Version, Commit, BuildDate = origVersion, origCommit, origBuildDate
	})

	Version, Commit, BuildDate = "", "", ""

	info := Get()
	if info.Version != "dev" || info.Commit != "dev" || info.BuildDate != "dev" {
		t.Fatalf("expected dev defaults, got %+v", info)
	}
}

func TestGetUsesOverrides(t *testing.T) {
	origVersion, origCommit, origBuildDate := Version, Commit, BuildDate
	t.Cleanup(func() {
		Version, Commit, BuildDate = origVersion, origCommit, origBuildDate
	})

	Version, Commit, BuildDate = "v1.2.3", "abc123", "2025-12-18"

	info := Get()
	if info.Version != "v1.2.3" || info.Commit != "abc123" || info.BuildDate != "2025-12-18" {
		t.Fatalf("unexpected overrides: %+v", info)
	}
}
