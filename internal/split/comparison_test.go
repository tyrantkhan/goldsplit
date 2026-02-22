package split

import "testing"

func TestComputeDelta(t *testing.T) {
	// Behind by 500ms.
	d := ComputeDelta(1500, 1000)
	if d != 500 {
		t.Fatalf("expected 500, got %d", d)
	}

	// Ahead by 200ms.
	d = ComputeDelta(800, 1000)
	if d != -200 {
		t.Fatalf("expected -200, got %d", d)
	}

	// Exact.
	d = ComputeDelta(1000, 1000)
	if d != 0 {
		t.Fatalf("expected 0, got %d", d)
	}
}

func TestComputeSplitDeltas(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})
	att.AddAttempt([]int64{1000, 2500, 4000}, true)

	// Faster run.
	deltas := ComputeSplitDeltas(att, []int64{900, 2300, 3800}, "personal_best")

	if len(deltas) != 3 {
		t.Fatalf("expected 3 deltas, got %d", len(deltas))
	}

	if !deltas[0].IsAhead || deltas[0].DeltaMS != -100 {
		t.Fatalf("delta[0] expected ahead by -100, got %+v", deltas[0])
	}

	if !deltas[0].GainedTime {
		t.Fatal("delta[0] expected gainedTime (900 < 1000)")
	}

	if !deltas[2].IsAhead || deltas[2].DeltaMS != -200 {
		t.Fatalf("delta[2] expected ahead by -200, got %+v", deltas[2])
	}

	// Seg C: time = 3800-2300 = 1500, comp seg = 4000-2500 = 1500 => not gained (equal).
	if deltas[2].GainedTime {
		t.Fatal("delta[2] expected no gainedTime (1500 == 1500)")
	}
}

func TestGainedTimeMixed(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})
	att.AddAttempt([]int64{1000, 2500, 4000}, true)

	// Seg A: faster (800 < 1000), Seg B: slower (1900 > 1500), Seg C: faster (1100 < 1500).
	// Cumulative: A=800(-200), B=2700(+200), C=3800(-200).
	deltas := ComputeSplitDeltas(att, []int64{800, 2700, 3800}, "personal_best")

	// Seg A: ahead + gained.
	if !deltas[0].IsAhead || !deltas[0].GainedTime {
		t.Fatalf("delta[0] expected ahead+gained, got %+v", deltas[0])
	}

	// Seg B: behind + lost (seg time 1900 > 1500).
	if deltas[1].IsAhead || deltas[1].GainedTime {
		t.Fatalf("delta[1] expected behind+lost, got %+v", deltas[1])
	}

	// Seg C: ahead + gained (seg time 1100 < 1500).
	if !deltas[2].IsAhead || !deltas[2].GainedTime {
		t.Fatalf("delta[2] expected ahead+gained, got %+v", deltas[2])
	}
}

func TestComputeSplitDeltasNoPB(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	deltas := ComputeSplitDeltas(att, []int64{1000, 2000}, "personal_best")

	// No PB, so no delta values.
	if deltas[0].DeltaMS != 0 || deltas[1].DeltaMS != 0 {
		t.Fatalf("expected zero deltas with no PB, got %+v", deltas)
	}
}

func TestComputeSplitDeltasSkipped(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})
	att.AddAttempt([]int64{1000, 2500, 4000}, true)

	deltas := ComputeSplitDeltas(att, []int64{0, 2300, 3800}, "personal_best")

	if !deltas[0].Skipped {
		t.Fatal("expected delta[0] to be skipped")
	}
}

func TestBestSegmentDetection(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	// Add a run to establish best segments.
	att.AddAttempt([]int64{1000, 2500}, true)

	// Run where segment 0 is better.
	deltas := ComputeSplitDeltas(att, []int64{800, 2400}, "personal_best")

	if !deltas[0].IsBestEver {
		t.Fatal("expected delta[0] to be best ever")
	}

	if deltas[1].IsBestEver {
		t.Fatal("expected delta[1] to not be best ever (1600 > 1500)")
	}
}

func TestPersonalBestSplitsIgnoresSkippedFinal(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	// Attempt with skipped final.
	att.AddAttempt([]int64{500, 0}, true)

	// Attempt with valid final.
	att.AddAttempt([]int64{1000, 2000}, true)

	pb := att.PersonalBestSplits()
	if pb == nil {
		t.Fatal("expected non-nil PB splits")
	}

	// Should pick the valid attempt, not the one with finalTime=0.
	if pb[1] != 2000 {
		t.Fatalf("expected PB final 2000, got %d", pb[1])
	}
}

func TestComputeSplitDeltasSkippedPBSegments(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C", "D"})

	// PB run with segment B skipped (0).
	att.AddAttempt([]int64{1000, 0, 3000, 5000}, true)

	// Current run with all segments.
	deltas := ComputeSplitDeltas(att, []int64{1100, 2000, 3200, 4800}, "personal_best")

	// Segment A: has PB data, should compute delta.
	if deltas[0].DeltaMS != 100 {
		t.Fatalf("delta[0] expected 100, got %d", deltas[0].DeltaMS)
	}

	// Segment B: PB is 0, should not compute delta.
	if deltas[1].DeltaMS != 0 {
		t.Fatalf("delta[1] expected 0 (no PB data), got %d", deltas[1].DeltaMS)
	}

	// Segment C: has PB data, should compute delta.
	if deltas[2].DeltaMS != 200 {
		t.Fatalf("delta[2] expected 200, got %d", deltas[2].DeltaMS)
	}
}

func TestBestSegmentsFromIncompleteRuns(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})

	// Incomplete run with only 2 splits.
	att.AddAttempt([]int64{800, 1800}, false)

	// Complete run.
	att.AddAttempt([]int64{1000, 2500, 4000}, true)

	best := att.BestSegments()

	// Segment A: 800 from incomplete run beats 1000 from complete run.
	if best[0] != 800 {
		t.Fatalf("expected best[0] 800, got %d", best[0])
	}

	// Segment B: 1000 from incomplete (1800-800) beats 1500 from complete (2500-1000).
	if best[1] != 1000 {
		t.Fatalf("expected best[1] 1000, got %d", best[1])
	}

	// Segment C: only from complete run (4000-2500 = 1500).
	if best[2] != 1500 {
		t.Fatalf("expected best[2] 1500, got %d", best[2])
	}
}

func TestComputeSplitDeltasEffectivelySkipped(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C", "D"})
	att.AddAttempt([]int64{1000, 2500, 4000, 5500}, true)

	// Current run where segments B and C are effectively skipped (same cumulative as A).
	deltas := ComputeSplitDeltas(att, []int64{1000, 1000, 1000, 2000}, "personal_best")

	// Segments B and C should be marked as skipped.
	if !deltas[1].Skipped {
		t.Fatal("expected delta[1] to be skipped (effectively skipped)")
	}

	if !deltas[2].Skipped {
		t.Fatal("expected delta[2] to be skipped (effectively skipped)")
	}

	// Segment D should not be skipped (2000-1000 = 1000 > 0).
	if deltas[3].Skipped {
		t.Fatal("expected delta[3] to not be skipped")
	}
}
