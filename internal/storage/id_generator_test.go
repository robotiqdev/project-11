package storage

import (
	"sync"
	"testing"
)

// TestNextID_StartsAtOne verifies that the first call to NextID() returns 1.
func TestNextID_StartsAtOne(t *testing.T) {
	g := NewIDGenerator()
	id := g.NextID()
	if id != 1 {
		t.Errorf("expected first ID to be 1, got %d", id)
	}
}

// TestNextID_IncrementsMonotonically verifies that successive calls return strictly increasing values.
func TestNextID_IncrementsMonotonically(t *testing.T) {
	g := NewIDGenerator()
	prev := g.NextID()
	for i := 0; i < 99; i++ {
		next := g.NextID()
		if next <= prev {
			t.Errorf("expected monotonically increasing IDs, got %d after %d", next, prev)
		}
		prev = next
	}
}

// TestNextID_SequentialValues verifies that IDs are exactly N+1 after N calls.
func TestNextID_SequentialValues(t *testing.T) {
	const n = 50
	g := NewIDGenerator()
	var last int64
	for i := 0; i < n; i++ {
		last = g.NextID()
	}
	if last != n {
		t.Errorf("after %d calls, expected last ID to be %d, got %d", n, n, last)
	}
}

// TestNextID_AfterNCalls_NextIsNPlusOne verifies that after N calls, the next value is N+1.
func TestNextID_AfterNCalls_NextIsNPlusOne(t *testing.T) {
	const n = 42
	g := NewIDGenerator()
	for i := 0; i < n; i++ {
		g.NextID()
	}
	next := g.NextID()
	if next != n+1 {
		t.Errorf("after %d calls, expected next ID to be %d, got %d", n, n+1, next)
	}
}

// TestNextID_NeverZero verifies that NextID() never returns zero.
func TestNextID_NeverZero(t *testing.T) {
	g := NewIDGenerator()
	for i := 0; i < 1000; i++ {
		id := g.NextID()
		if id == 0 {
			t.Errorf("NextID() returned zero at call %d", i+1)
		}
	}
}

// TestNextID_NeverNegative verifies that NextID() never returns a negative value.
func TestNextID_NeverNegative(t *testing.T) {
	g := NewIDGenerator()
	for i := 0; i < 1000; i++ {
		id := g.NextID()
		if id < 0 {
			t.Errorf("NextID() returned negative value %d at call %d", id, i+1)
		}
	}
}

// TestNextID_ConcurrentUnique verifies that 100 goroutines each calling NextID() 100 times
// produce 10000 unique values with no gaps.
func TestNextID_ConcurrentUnique(t *testing.T) {
	const goroutines = 100
	const callsPerGoroutine = 100
	const total = goroutines * callsPerGoroutine

	g := NewIDGenerator()

	ids := make(chan int64, total)
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				ids <- g.NextID()
			}
		}()
	}

	wg.Wait()
	close(ids)

	seen := make(map[int64]bool, total)
	for id := range ids {
		if seen[id] {
			t.Errorf("duplicate ID: %d", id)
		}
		seen[id] = true
	}

	if len(seen) != total {
		t.Errorf("expected %d unique IDs, got %d", total, len(seen))
	}

	// Verify no gaps: all values from 1 to total must be present.
	for i := int64(1); i <= total; i++ {
		if !seen[i] {
			t.Errorf("gap in IDs: missing %d", i)
		}
	}
}

// TestNextID_ConcurrentNoGaps verifies sequential IDs under concurrent load have no missing values.
func TestNextID_ConcurrentNoGaps(t *testing.T) {
	const goroutines = 10
	const callsPerGoroutine = 100
	const total = goroutines * callsPerGoroutine

	g := NewIDGenerator()
	results := make([]int64, total)
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				id := g.NextID()
				mu.Lock()
				results = append(results[:0:0], results...) // no-op to force usage
				_ = id
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// The counter must have reached exactly `total` after all calls.
	finalID := g.NextID()
	if finalID != total+1 {
		t.Errorf("after %d concurrent calls, expected next ID to be %d, got %d", total, total+1, finalID)
	}
}

// TestNewIDGenerator_IndependentInstances verifies that separate IDGenerator instances
// have independent counters.
func TestNewIDGenerator_IndependentInstances(t *testing.T) {
	g1 := NewIDGenerator()
	g2 := NewIDGenerator()

	id1 := g1.NextID()
	id2 := g2.NextID()

	if id1 != 1 {
		t.Errorf("g1: expected first ID 1, got %d", id1)
	}
	if id2 != 1 {
		t.Errorf("g2: expected first ID 1, got %d", id2)
	}

	// Advance g1 further and confirm g2 is unaffected.
	g1.NextID()
	g1.NextID()
	id2next := g2.NextID()
	if id2next != 2 {
		t.Errorf("g2 should not be affected by g1; expected 2, got %d", id2next)
	}
}
