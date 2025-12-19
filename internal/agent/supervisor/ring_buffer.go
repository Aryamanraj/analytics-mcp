package supervisor

import "sync"

type ringBuffer struct {
	mu    sync.Mutex
	lines []string
	next  int
	count int
}

func newRingBuffer(size int) *ringBuffer {
	if size <= 0 {
		size = 200
	}
	return &ringBuffer{lines: make([]string, size)}
}

func (r *ringBuffer) Add(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lines[r.next] = line
	r.next = (r.next + 1) % len(r.lines)
	if r.count < len(r.lines) {
		r.count++
	}
}

func (r *ringBuffer) Tail(n int) []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if n <= 0 || r.count == 0 {
		return nil
	}

	if n > r.count {
		n = r.count
	}

	start := (r.next - n + len(r.lines)) % len(r.lines)
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = r.lines[(start+i)%len(r.lines)]
	}

	return out
}
