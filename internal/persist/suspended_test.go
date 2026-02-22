package persist

import (
	"testing"
)

func TestSuspendedRunRoundTrip(t *testing.T) {
	store := tempStore(t)

	run := &SuspendedRun{
		TemplateID:     "tmpl-1",
		AttemptsID:     "att-1",
		ElapsedMS:      12345,
		CurrentSegment: 2,
		SplitTimesMS:   []int64{1000, 3000},
		SegmentTimesMS: []int64{1000, 2000},
		SuspendedAt:    1700000000,
	}

	if err := store.SaveSuspendedRun(run); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := store.LoadSuspendedRun()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded == nil {
		t.Fatal("expected non-nil suspended run")

		return
	}

	if loaded.TemplateID != run.TemplateID {
		t.Fatalf("templateId: got %q, want %q", loaded.TemplateID, run.TemplateID)
	}

	if loaded.AttemptsID != run.AttemptsID {
		t.Fatalf("attemptsId: got %q, want %q", loaded.AttemptsID, run.AttemptsID)
	}

	if loaded.ElapsedMS != run.ElapsedMS {
		t.Fatalf("elapsedMs: got %d, want %d", loaded.ElapsedMS, run.ElapsedMS)
	}

	if loaded.CurrentSegment != run.CurrentSegment {
		t.Fatalf("currentSegment: got %d, want %d", loaded.CurrentSegment, run.CurrentSegment)
	}

	if len(loaded.SplitTimesMS) != len(run.SplitTimesMS) {
		t.Fatalf("splitTimesMs length: got %d, want %d", len(loaded.SplitTimesMS), len(run.SplitTimesMS))
	}

	for i, v := range loaded.SplitTimesMS {
		if v != run.SplitTimesMS[i] {
			t.Fatalf("splitTimesMs[%d]: got %d, want %d", i, v, run.SplitTimesMS[i])
		}
	}

	if loaded.SuspendedAt != run.SuspendedAt {
		t.Fatalf("suspendedAt: got %d, want %d", loaded.SuspendedAt, run.SuspendedAt)
	}
}

func TestLoadSuspendedRunMissing(t *testing.T) {
	store := tempStore(t)

	loaded, err := store.LoadSuspendedRun()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if loaded != nil {
		t.Fatal("expected nil when file does not exist")
	}
}

func TestDeleteSuspendedRun(t *testing.T) {
	store := tempStore(t)

	run := &SuspendedRun{
		TemplateID: "tmpl-1",
		AttemptsID: "att-1",
		ElapsedMS:  5000,
	}

	if err := store.SaveSuspendedRun(run); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if err := store.DeleteSuspendedRun(); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	loaded, err := store.LoadSuspendedRun()
	if err != nil {
		t.Fatalf("load after delete failed: %v", err)
	}

	if loaded != nil {
		t.Fatal("expected nil after delete")
	}
}

func TestDeleteSuspendedRunIdempotent(t *testing.T) {
	store := tempStore(t)

	// Deleting when no file exists should not error.
	if err := store.DeleteSuspendedRun(); err != nil {
		t.Fatalf("first delete failed: %v", err)
	}

	if err := store.DeleteSuspendedRun(); err != nil {
		t.Fatalf("second delete failed: %v", err)
	}
}
