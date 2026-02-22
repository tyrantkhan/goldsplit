package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func segments() []string {
	return []string{"Segment 1", "Segment 2", "Segment 3"}
}

func TestNewEngine(t *testing.T) {
	e := New(segments(), nil, nil)
	if e.CurrentState() != Idle {
		t.Fatalf("expected Idle, got %v", e.CurrentState())
	}
}

func TestStartTransition(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()

	if e.CurrentState() != Running {
		t.Fatalf("expected Running, got %v", e.CurrentState())
	}

	e.Reset()
}

func TestStartOnlyFromIdle(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()
	e.Pause()

	// Start should be a no-op from Paused.
	e.Start()

	if e.CurrentState() != Paused {
		t.Fatalf("expected Paused, got %v", e.CurrentState())
	}

	e.Reset()
}

func TestPauseResume(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()
	time.Sleep(50 * time.Millisecond)
	e.Pause()

	if e.CurrentState() != Paused {
		t.Fatalf("expected Paused, got %v", e.CurrentState())
	}

	pausedElapsed := e.ElapsedMS()
	time.Sleep(50 * time.Millisecond)

	// Elapsed should not change while paused.
	if e.ElapsedMS() != pausedElapsed {
		t.Fatalf("elapsed changed during pause: %d -> %d", pausedElapsed, e.ElapsedMS())
	}

	e.Resume()

	if e.CurrentState() != Running {
		t.Fatalf("expected Running after resume, got %v", e.CurrentState())
	}

	time.Sleep(50 * time.Millisecond)

	if e.ElapsedMS() <= pausedElapsed {
		t.Fatalf("elapsed should have increased after resume")
	}

	e.Reset()
}

func TestSplitCompleteRun(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()
	time.Sleep(20 * time.Millisecond)

	e.Split() // Segment 1
	e.Split() // Segment 2
	e.Split() // Segment 3 -> Finished

	if e.CurrentState() != Finished {
		t.Fatalf("expected Finished after all splits, got %v", e.CurrentState())
	}

	splits := e.SplitTimesMS()
	if len(splits) != 3 {
		t.Fatalf("expected 3 splits, got %d", len(splits))
	}

	// Splits should be monotonically increasing.
	for i := 1; i < len(splits); i++ {
		if splits[i] < splits[i-1] {
			t.Fatalf("split %d (%d) < split %d (%d)", i, splits[i], i-1, splits[i-1])
		}
	}
}

func TestUndoSplit(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()
	time.Sleep(20 * time.Millisecond)

	e.Split()

	if e.CurrentSegment() != 1 {
		t.Fatalf("expected segment 1, got %d", e.CurrentSegment())
	}

	e.UndoSplit()

	if e.CurrentSegment() != 0 {
		t.Fatalf("expected segment 0 after undo, got %d", e.CurrentSegment())
	}

	e.Reset()
}

func TestUndoSplitAtZero(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()

	// Should be a no-op.
	e.UndoSplit()

	if e.CurrentSegment() != 0 {
		t.Fatalf("expected segment 0, got %d", e.CurrentSegment())
	}

	e.Reset()
}

func TestSkipSplit(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()
	time.Sleep(20 * time.Millisecond)

	e.SkipSplit() // Skip segment 1
	e.Split()     // Segment 2
	e.Split()     // Segment 3 -> Finished

	if e.CurrentState() != Finished {
		t.Fatalf("expected Finished, got %v", e.CurrentState())
	}

	splits := e.SplitTimesMS()
	if splits[0] != 0 {
		t.Fatalf("expected skipped split to be 0, got %d", splits[0])
	}
}

func TestResetFromRunning(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()
	time.Sleep(20 * time.Millisecond)
	e.Reset()

	if e.CurrentState() != Idle {
		t.Fatalf("expected Idle after reset, got %v", e.CurrentState())
	}
}

func TestResetFromFinished(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()

	e.Split()
	e.Split()
	e.Split()

	e.Reset()

	if e.CurrentState() != Idle {
		t.Fatalf("expected Idle after reset, got %v", e.CurrentState())
	}
}

func TestResetFromIdleNoop(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Reset() // Should not panic or change state.

	if e.CurrentState() != Idle {
		t.Fatalf("expected Idle, got %v", e.CurrentState())
	}
}

func TestTickCallback(t *testing.T) {
	var tickCount atomic.Int64

	e := New(segments(), func(_ TickData) {
		tickCount.Add(1)
	}, nil)

	e.Start()
	time.Sleep(100 * time.Millisecond)
	e.Reset()

	// At 15ms intervals over 100ms, we expect roughly 5-7 ticks.
	count := tickCount.Load()
	if count < 3 || count > 15 {
		t.Fatalf("unexpected tick count: %d", count)
	}
}

func TestStateChangeCallback(t *testing.T) {
	var states []State

	e := New(segments(), nil, func(s State) {
		states = append(states, s)
	})

	e.Start()
	e.Pause()
	e.Resume()
	e.Reset()

	expected := []State{Running, Paused, Running, Idle}
	if len(states) != len(expected) {
		t.Fatalf("expected %d state changes, got %d: %v", len(expected), len(states), states)
	}

	for i, s := range states {
		if s != expected[i] {
			t.Fatalf("state[%d]: expected %v, got %v", i, expected[i], s)
		}
	}
}

func TestElapsedAccuracy(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()
	time.Sleep(100 * time.Millisecond)
	elapsed := e.ElapsedMS()
	e.Reset()

	// Should be roughly 100ms, allow 50ms-200ms range for CI.
	if elapsed < 50 || elapsed > 200 {
		t.Fatalf("elapsed %dms out of expected range [50, 200]", elapsed)
	}
}

func TestPauseAccumulation(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()
	time.Sleep(50 * time.Millisecond)
	e.Pause()
	time.Sleep(100 * time.Millisecond) // This should not count.
	e.Resume()
	time.Sleep(50 * time.Millisecond)
	elapsed := e.ElapsedMS()
	e.Reset()

	// ~100ms of actual running, 100ms paused should not count.
	// Allow 50ms-200ms.
	if elapsed < 50 || elapsed > 200 {
		t.Fatalf("elapsed %dms out of expected range [50, 200] (pause should not count)", elapsed)
	}
}

func TestSetSegments(t *testing.T) {
	e := New(segments(), nil, nil)
	newSegs := []string{"A", "B"}
	e.SetSegments(newSegs)

	e.Start()
	e.Split()
	e.Split()

	if e.CurrentState() != Finished {
		t.Fatalf("expected Finished with 2 segments, got %v", e.CurrentState())
	}

	e.Reset()
}

func TestRestoreSetssPausedWithCorrectElapsed(t *testing.T) {
	e := New(segments(), nil, nil)

	splits := []int64{1000, 3000}
	segs := []int64{1000, 2000}
	e.Restore(5000, 2, splits, segs)

	if e.CurrentState() != Paused {
		t.Fatalf("expected Paused, got %v", e.CurrentState())
	}

	elapsed := e.ElapsedMS()
	if elapsed < 4900 || elapsed > 5100 {
		t.Fatalf("expected ~5000ms elapsed, got %d", elapsed)
	}

	if e.CurrentSegment() != 2 {
		t.Fatalf("expected segment 2, got %d", e.CurrentSegment())
	}

	gotSplits := e.SplitTimesMS()
	if len(gotSplits) != 2 || gotSplits[0] != 1000 || gotSplits[1] != 3000 {
		t.Fatalf("unexpected splits: %v", gotSplits)
	}

	gotSegs := e.SegmentTimesMS()
	if len(gotSegs) != 2 || gotSegs[0] != 1000 || gotSegs[1] != 2000 {
		t.Fatalf("unexpected segments: %v", gotSegs)
	}

	e.Reset()
}

func TestRestoreOnlyFromIdle(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()

	// Restore should be a no-op when not Idle.
	e.Restore(5000, 2, []int64{1000, 3000}, []int64{1000, 2000})

	if e.CurrentState() != Running {
		t.Fatalf("expected Running, got %v", e.CurrentState())
	}

	e.Reset()
}

func TestRestoreThenResumeElapsedIncreases(t *testing.T) {
	e := New(segments(), nil, nil)

	e.Restore(5000, 1, []int64{1000}, []int64{1000})

	pausedElapsed := e.ElapsedMS()

	e.Resume()

	if e.CurrentState() != Running {
		t.Fatalf("expected Running after resume, got %v", e.CurrentState())
	}

	time.Sleep(50 * time.Millisecond)

	if e.ElapsedMS() <= pausedElapsed {
		t.Fatal("elapsed should have increased after resume")
	}

	e.Reset()
}

func TestSetSegmentsWhileRunning(t *testing.T) {
	e := New(segments(), nil, nil)
	e.Start()
	e.SetSegments([]string{"A"}) // Should be a no-op.

	e.Split()

	if e.CurrentState() == Finished {
		t.Fatal("SetSegments should have been ignored while running")
	}

	e.Reset()
}
