package supervisor

import (
	"context"
	"testing"
	"time"
)

func TestRestartBackoffWithFailingCommand(t *testing.T) {
	cfg := Config{
		ChatPath:         "/bin/sh",
		ChatArgs:         []string{"-c", "exit 1"},
		MCPPath:          "/bin/sh",
		MCPArgs:          []string{"-c", "exit 1"},
		BufferLines:      20,
		InitialBackoff:   20 * time.Millisecond,
		MaxBackoff:       50 * time.Millisecond,
		TerminateTimeout: 200 * time.Millisecond,
	}

	sup := New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	if err := sup.Start(ctx); err != nil {
		t.Fatalf("failed to start supervisor: %v", err)
	}

	sup.Wait()

	status := sup.Status()
	if len(status.Components) != 2 {
		t.Fatalf("expected 2 components, got %d", len(status.Components))
	}

	for _, comp := range status.Components {
		if comp.Restarts < 1 {
			t.Fatalf("component %s should have restarted at least once", comp.Name)
		}
		if comp.Restarts > 5 {
			t.Fatalf("component %s restart count too high without backoff cap: %d", comp.Name, comp.Restarts)
		}
	}
}
