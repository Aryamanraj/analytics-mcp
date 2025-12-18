package update

import "testing"

func TestStatusPersistence(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PAYRAM_AGENT_HOME", home)

	st := UpdateStatus{CurrentVersion: "1.0.0", PreviousVersion: "0.9.0"}
	st.MarkSuccess("1.0.0", "0.9.0")
	st.MarkFailure("ERR", "boom")
	st.MarkAttempt()

	if err := SaveStatus(st); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := LoadStatus()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if loaded.CurrentVersion != st.CurrentVersion {
		t.Fatalf("current mismatch")
	}
	if loaded.PreviousVersion != st.PreviousVersion {
		t.Fatalf("previous mismatch")
	}
	if loaded.LastAttemptAt.IsZero() {
		t.Fatalf("expected attempt timestamp")
	}
	if loaded.LastAttemptVersion != "" {
		t.Fatalf("expected empty attempt version")
	}
	if !loaded.InProgress {
		t.Fatalf("expected in progress true")
	}
}
