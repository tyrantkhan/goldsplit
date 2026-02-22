package persist

import (
	"os"
	"testing"

	"goldsplit/internal/split"
)

func tempStore(t *testing.T) *Store {
	t.Helper()

	dir := t.TempDir()

	store, err := NewStore(dir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	return store
}

func TestSaveAndLoadTemplate(t *testing.T) {
	store := tempStore(t)
	tmpl := split.NewTemplate("t-1", "Super Mario 64", []string{"Bob-omb", "Whomp", "Jolly Roger"})

	if err := store.SaveTemplate(tmpl); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := store.LoadTemplate("t-1")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Name != "Super Mario 64" {
		t.Fatalf("expected name Super Mario 64, got %s", loaded.Name)
	}

	if len(loaded.SegmentNames) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(loaded.SegmentNames))
	}
}

func TestListTemplates(t *testing.T) {
	store := tempStore(t)

	t1 := split.NewTemplate("t-1", "Game A", []string{"Seg1"})
	t2 := split.NewTemplate("t-2", "Game B", []string{"Seg1", "Seg2"})

	if err := store.SaveTemplate(t1); err != nil {
		t.Fatalf("save t1 failed: %v", err)
	}

	if err := store.SaveTemplate(t2); err != nil {
		t.Fatalf("save t2 failed: %v", err)
	}

	summaries, err := store.ListTemplates()
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}
}

func TestDeleteTemplate(t *testing.T) {
	store := tempStore(t)

	tmpl := split.NewTemplate("del-1", "Game", []string{"Seg"})
	if err := store.SaveTemplate(tmpl); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if err := store.DeleteTemplate("del-1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err := store.LoadTemplate("del-1")
	if err == nil {
		t.Fatal("expected error loading deleted template")
	}
}

func TestDeleteTemplateCascadesAttempts(t *testing.T) {
	store := tempStore(t)

	tmpl := split.NewTemplate("t-1", "Game", []string{"Seg"})
	if err := store.SaveTemplate(tmpl); err != nil {
		t.Fatalf("save template failed: %v", err)
	}

	att := split.NewAttempts("a-1", "t-1", "", "Any%", []string{"Seg"})
	if err := store.SaveAttempts(att); err != nil {
		t.Fatalf("save attempts failed: %v", err)
	}

	if err := store.DeleteTemplate("t-1"); err != nil {
		t.Fatalf("delete template failed: %v", err)
	}

	_, err := store.LoadAttempts("a-1")
	if err == nil {
		t.Fatal("expected attempts to be deleted with template")
	}
}

func TestSaveAndLoadAttempts(t *testing.T) {
	store := tempStore(t)
	att := split.NewAttempts("a-1", "t-1", "", "Any%", []string{"A", "B"})
	att.AddAttempt([]int64{1000, 2000}, true)

	if err := store.SaveAttempts(att); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := store.LoadAttempts("a-1")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.CategoryName != "Any%" {
		t.Fatalf("expected category Any%%, got %s", loaded.CategoryName)
	}

	if len(loaded.Segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(loaded.Segments))
	}

	if len(loaded.History) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(loaded.History))
	}
}

func TestListAttemptsForTemplate(t *testing.T) {
	store := tempStore(t)

	a1 := split.NewAttempts("a-1", "t-1", "", "Any%", []string{"Seg"})
	a2 := split.NewAttempts("a-2", "t-1", "", "100%", []string{"Seg"})
	a3 := split.NewAttempts("a-3", "t-2", "", "Any%", []string{"Seg"})

	for _, a := range []*split.Attempts{a1, a2, a3} {
		if err := store.SaveAttempts(a); err != nil {
			t.Fatalf("save failed: %v", err)
		}
	}

	summaries, err := store.ListAttemptsForTemplate("t-1")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries for t-1, got %d", len(summaries))
	}
}

func TestDeleteAttempts(t *testing.T) {
	store := tempStore(t)

	att := split.NewAttempts("del-1", "t-1", "", "Any%", []string{"Seg"})
	if err := store.SaveAttempts(att); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if err := store.DeleteAttempts("del-1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err := store.LoadAttempts("del-1")
	if err == nil {
		t.Fatal("expected error loading deleted attempts")
	}
}

func TestLoadNonexistentTemplate(t *testing.T) {
	store := tempStore(t)

	_, err := store.LoadTemplate("nonexistent")
	if err == nil {
		t.Fatal("expected error loading nonexistent template")
	}
}

func TestLoadNonexistentAttempts(t *testing.T) {
	store := tempStore(t)

	_, err := store.LoadAttempts("nonexistent")
	if err == nil {
		t.Fatal("expected error loading nonexistent attempts")
	}
}

func TestAtomicWrite(t *testing.T) {
	store := tempStore(t)

	tmpl := split.NewTemplate("atomic-1", "Game", []string{"Seg"})
	if err := store.SaveTemplate(tmpl); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	entries, err := os.ReadDir(store.baseDir + "/templates")
	if err != nil {
		t.Fatalf("readdir failed: %v", err)
	}

	for _, e := range entries {
		if e.Name()[len(e.Name())-4:] == ".tmp" {
			t.Fatalf("tmp file left behind: %s", e.Name())
		}
	}
}

func TestSaveUpdatesRoundTrip(t *testing.T) {
	store := tempStore(t)

	att := split.NewAttempts("rt-1", "t-1", "", "Any%", []string{"A", "B"})
	if err := store.SaveAttempts(att); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	att.AddAttempt([]int64{1000, 2000}, true)

	if err := store.SaveAttempts(att); err != nil {
		t.Fatalf("save updated failed: %v", err)
	}

	loaded, err := store.LoadAttempts("rt-1")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.AttemptCount != 1 {
		t.Fatalf("expected 1 attempt, got %d", loaded.AttemptCount)
	}

	pb := loaded.PersonalBestSplits()
	if pb == nil || pb[0] != 1000 {
		t.Fatalf("expected PB 1000, got %v", pb)
	}
}
