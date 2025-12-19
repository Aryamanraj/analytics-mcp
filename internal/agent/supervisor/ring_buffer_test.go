package supervisor

import "testing"

func TestRingBufferTail(t *testing.T) {
	buf := newRingBuffer(3)

	buf.Add("a")
	buf.Add("b")
	if got := buf.Tail(2); len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("unexpected tail: %v", got)
	}

	buf.Add("c")
	buf.Add("d") // overwrites "a"

	got := buf.Tail(3)
	expected := []string{"b", "c", "d"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(got))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("unexpected line %d: want %q got %q", i, expected[i], got[i])
		}
	}

	if empty := buf.Tail(0); empty != nil {
		t.Fatalf("expected nil for zero tail, got %v", empty)
	}
}
