package split

import "testing"

func TestNewAttempts(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"Seg1", "Seg2"})

	if att.ID != "a-1" {
		t.Fatalf("expected ID a-1, got %s", att.ID)
	}

	if att.TemplateID != "t-1" {
		t.Fatalf("expected TemplateID t-1, got %s", att.TemplateID)
	}

	if len(att.Segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(att.Segments))
	}

	if att.Segments[0].Name != "Seg1" {
		t.Fatalf("expected Seg1, got %s", att.Segments[0].Name)
	}
}

func TestAttemptsSegmentNames(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})
	names := att.SegmentNames()

	if len(names) != 3 || names[0] != "A" || names[1] != "B" || names[2] != "C" {
		t.Fatalf("unexpected segment names: %v", names)
	}
}

func TestAttemptsAddAttempt(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})
	att.AddAttempt([]int64{1000, 2000}, true)

	if att.AttemptCount != 1 {
		t.Fatalf("expected attempt count 1, got %d", att.AttemptCount)
	}

	if len(att.History) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(att.History))
	}

	if !att.History[0].Completed {
		t.Fatal("expected attempt to be completed")
	}
}

func TestAttemptsPersonalBestSplits(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	if att.PersonalBestSplits() != nil {
		t.Fatal("expected nil PB with no attempts")
	}

	att.AddAttempt([]int64{1000}, false)

	if att.PersonalBestSplits() != nil {
		t.Fatal("expected nil PB with no completed attempts")
	}

	att.AddAttempt([]int64{1000, 3000}, true)
	pb := att.PersonalBestSplits()

	if pb == nil || pb[1] != 3000 {
		t.Fatalf("expected PB [1000, 3000], got %v", pb)
	}

	att.AddAttempt([]int64{900, 2500}, true)
	pb = att.PersonalBestSplits()

	if pb[1] != 2500 {
		t.Fatalf("expected PB final 2500, got %d", pb[1])
	}

	att.AddAttempt([]int64{1100, 4000}, true)
	pb = att.PersonalBestSplits()

	if pb[1] != 2500 {
		t.Fatalf("expected PB final still 2500, got %d", pb[1])
	}
}

func TestAttemptsBestSegmentsCumulative(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})

	if att.BestSegmentsCumulative() != nil {
		t.Fatal("expected nil with no best segments")
	}

	att.Segments[0].BestSegmentMS = 500
	att.Segments[1].BestSegmentMS = 700
	att.Segments[2].BestSegmentMS = 800

	cum := att.BestSegmentsCumulative()
	if cum == nil {
		t.Fatal("expected non-nil cumulative splits")
	}

	if cum[0] != 500 || cum[1] != 1200 || cum[2] != 2000 {
		t.Fatalf("expected [500, 1200, 2000], got %v", cum)
	}
}

func TestAttemptsBestSegmentsCumulativePartial(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})
	att.Segments[0].BestSegmentMS = 500

	if att.BestSegmentsCumulative() != nil {
		t.Fatal("expected nil when a segment has no best")
	}
}

func TestAttemptsAverageSplits(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	if att.AverageSplits() != nil {
		t.Fatal("expected nil with no attempts")
	}

	att.AddAttempt([]int64{1000, 3000}, true)
	att.AddAttempt([]int64{2000, 5000}, true)

	avg := att.AverageSplits()
	if avg == nil {
		t.Fatal("expected non-nil averages")
	}

	if avg[0] != 1500 || avg[1] != 4000 {
		t.Fatalf("expected [1500, 4000], got %v", avg)
	}
}

func TestAttemptsAverageSplitsSkipsInvalid(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	att.AddAttempt([]int64{500}, false)
	att.AddAttempt([]int64{600, 0}, true)
	att.AddAttempt([]int64{1000, 2000}, true)

	avg := att.AverageSplits()
	if avg == nil {
		t.Fatal("expected non-nil averages")
	}

	if avg[0] != 1000 || avg[1] != 2000 {
		t.Fatalf("expected [1000, 2000], got %v", avg)
	}
}

func TestAttemptsLatestCompletedSplits(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	if att.LatestCompletedSplits() != nil {
		t.Fatal("expected nil with no attempts")
	}

	att.AddAttempt([]int64{1000, 3000}, true)
	att.AddAttempt([]int64{900, 2500}, true)

	latest := att.LatestCompletedSplits()
	if latest == nil {
		t.Fatal("expected non-nil latest splits")
	}

	if latest[0] != 900 || latest[1] != 2500 {
		t.Fatalf("expected [900, 2500], got %v", latest)
	}
}

func TestAttemptsLatestCompletedSplitsSkipsInvalid(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})

	att.AddAttempt([]int64{1000, 3000}, true)
	att.AddAttempt([]int64{500}, false)
	att.AddAttempt([]int64{800, 0}, true)

	latest := att.LatestCompletedSplits()
	if latest == nil {
		t.Fatal("expected non-nil")
	}

	if latest[0] != 1000 || latest[1] != 3000 {
		t.Fatalf("expected [1000, 3000], got %v", latest)
	}
}

func TestDeleteAttempt(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})
	att.AddAttempt([]int64{1000, 3000}, true)
	UpdatePersonalBest(att, []int64{1000, 3000})
	att.AddAttempt([]int64{900, 2500}, true)
	UpdatePersonalBest(att, []int64{900, 2500})

	// Delete the PB attempt.
	if !att.DeleteAttempt(2) {
		t.Fatal("expected successful delete")
	}

	if len(att.History) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(att.History))
	}

	// PB should have been recalculated to the remaining attempt.
	if att.Segments[1].PersonalBestMS != 3000 {
		t.Fatalf("expected recalculated PB 3000, got %d", att.Segments[1].PersonalBestMS)
	}

	// Delete nonexistent.
	if att.DeleteAttempt(99) {
		t.Fatal("expected false for nonexistent attempt")
	}
}

func TestEditAttemptSplits(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})
	att.AddAttempt([]int64{1000, 3000}, true)
	UpdatePersonalBest(att, []int64{1000, 3000})

	// Edit to a faster time.
	if !att.EditAttemptSplits(1, []int64{800, 2000}) {
		t.Fatal("expected successful edit")
	}

	if att.History[0].SplitTimesMS[1] != 2000 {
		t.Fatalf("expected edited split 2000, got %d", att.History[0].SplitTimesMS[1])
	}

	// PB should reflect the new times.
	if att.Segments[1].PersonalBestMS != 2000 {
		t.Fatalf("expected PB 2000, got %d", att.Segments[1].PersonalBestMS)
	}

	// Edit nonexistent.
	if att.EditAttemptSplits(99, []int64{100, 200}) {
		t.Fatal("expected false for nonexistent attempt")
	}
}

func TestEditAttemptSplitsRejectsMismatchedLength(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})
	att.AddAttempt([]int64{1000, 3000}, true)

	if att.EditAttemptSplits(1, []int64{1000}) {
		t.Fatal("expected false for mismatched length (too few)")
	}

	if att.EditAttemptSplits(1, []int64{1000, 2000, 3000}) {
		t.Fatal("expected false for mismatched length (too many)")
	}
}

func TestEditAttemptSplitsRejectsNonMonotonic(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})
	att.AddAttempt([]int64{1000, 2000, 3000}, true)

	// Second split smaller than first.
	if att.EditAttemptSplits(1, []int64{2000, 1000, 3000}) {
		t.Fatal("expected false for non-monotonic splits")
	}

	// Equal values (not strictly increasing).
	if att.EditAttemptSplits(1, []int64{1000, 1000, 3000}) {
		t.Fatal("expected false for equal consecutive splits")
	}
}

func TestEditAttemptSplitsAcceptsSkips(t *testing.T) {
	att := NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B", "C"})
	att.AddAttempt([]int64{1000, 2000, 3000}, true)

	// Zero (skip) in the middle is valid.
	if !att.EditAttemptSplits(1, []int64{1000, 0, 3000}) {
		t.Fatal("expected true for valid edit with skip")
	}

	if att.History[0].SplitTimesMS[1] != 0 {
		t.Fatalf("expected 0, got %d", att.History[0].SplitTimesMS[1])
	}
}
