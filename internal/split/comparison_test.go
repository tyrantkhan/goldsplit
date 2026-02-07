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

func TestUpdatePersonalBest(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})

	// First completed run sets PB.
	splits := []int64{1000, 2500, 4000}
	att.AddAttempt(splits, true)
	UpdatePersonalBest(att, splits)

	if att.Segments[0].PersonalBestMS != 1000 {
		t.Fatalf("expected PB[0] 1000, got %d", att.Segments[0].PersonalBestMS)
	}

	if att.Segments[0].BestSegmentMS != 1000 {
		t.Fatalf("expected best seg[0] 1000, got %d", att.Segments[0].BestSegmentMS)
	}

	// Faster run updates PB.
	splits2 := []int64{900, 2300, 3800}
	att.AddAttempt(splits2, true)
	UpdatePersonalBest(att, splits2)

	if att.Segments[0].PersonalBestMS != 900 {
		t.Fatalf("expected PB[0] 900, got %d", att.Segments[0].PersonalBestMS)
	}

	// Slower overall but with a gold split on segment 1.
	splits3 := []int64{1100, 2100, 4100}
	att.AddAttempt(splits3, true)
	UpdatePersonalBest(att, splits3)

	// PB should not change (4100 > 3800).
	if att.Segments[2].PersonalBestMS != 3800 {
		t.Fatalf("expected PB[2] still 3800, got %d", att.Segments[2].PersonalBestMS)
	}

	// But best segment for seg 1 should update (1000ms vs 1400ms from PB).
	// Seg 1 time: 2100 - 1100 = 1000, previous best was 2300-900 = 1400.
	if att.Segments[1].BestSegmentMS != 1000 {
		t.Fatalf("expected best seg[1] 1000, got %d", att.Segments[1].BestSegmentMS)
	}
}

func TestUpdatePersonalBestWithSkips(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})
	splits := []int64{1000, 0}
	UpdatePersonalBest(att, splits)

	// Skipped final segment: no PB should be set.
	if att.Segments[0].PersonalBestMS != 0 {
		t.Fatalf("expected PB[0] 0 (skipped final), got %d", att.Segments[0].PersonalBestMS)
	}

	if att.Segments[1].BestSegmentMS != 0 {
		t.Fatalf("expected best seg[1] 0, got %d", att.Segments[1].BestSegmentMS)
	}

	// But best segment for seg 0 should still update.
	if att.Segments[0].BestSegmentMS != 1000 {
		t.Fatalf("expected best seg[0] 1000, got %d", att.Segments[0].BestSegmentMS)
	}
}

func TestSkipDoesNotWipePB(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})

	// First run sets PB.
	splits := []int64{1000, 2500, 4000}
	att.AddAttempt(splits, true)
	UpdatePersonalBest(att, splits)

	if att.Segments[1].PersonalBestMS != 2500 {
		t.Fatalf("expected PB[1] 2500, got %d", att.Segments[1].PersonalBestMS)
	}

	// Second run with middle segment skipped but faster final time.
	splits2 := []int64{800, 0, 3500}
	att.AddAttempt(splits2, true)
	UpdatePersonalBest(att, splits2)

	// New PB (3500 < 4000), but skipped segment should preserve old PB value.
	if att.Segments[0].PersonalBestMS != 800 {
		t.Fatalf("expected PB[0] 800, got %d", att.Segments[0].PersonalBestMS)
	}

	if att.Segments[1].PersonalBestMS != 2500 {
		t.Fatalf("expected PB[1] preserved at 2500, got %d", att.Segments[1].PersonalBestMS)
	}

	if att.Segments[2].PersonalBestMS != 3500 {
		t.Fatalf("expected PB[2] 3500, got %d", att.Segments[2].PersonalBestMS)
	}
}

func TestSkipFinalSegmentDoesNotSetPB(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	// First run sets PB.
	splits := []int64{1000, 2000}
	att.AddAttempt(splits, true)
	UpdatePersonalBest(att, splits)

	// Run with final segment skipped should not change PB.
	splits2 := []int64{500, 0}
	att.AddAttempt(splits2, true)
	UpdatePersonalBest(att, splits2)

	if att.Segments[0].PersonalBestMS != 1000 {
		t.Fatalf("expected PB[0] preserved at 1000, got %d", att.Segments[0].PersonalBestMS)
	}

	if att.Segments[1].PersonalBestMS != 2000 {
		t.Fatalf("expected PB[1] preserved at 2000, got %d", att.Segments[1].PersonalBestMS)
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

func TestBestSegmentDetection(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	// Set up existing best segments.
	att.Segments[0].BestSegmentMS = 1000
	att.Segments[1].BestSegmentMS = 1500

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
