package persist

import (
	"os"
	"testing"
)

func TestDefaultSettingsOnMissingFile(t *testing.T) {
	store := tempStore(t)

	settings, err := store.LoadSettings()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	defaults := DefaultSettings()
	if settings.Hotkeys.StartSplit != defaults.Hotkeys.StartSplit {
		t.Fatalf("expected default startSplit %q, got %q", defaults.Hotkeys.StartSplit, settings.Hotkeys.StartSplit)
	}

	if settings.AlwaysOnTop != false {
		t.Fatal("expected alwaysOnTop false")
	}

	if settings.Comparison != "personal_best" {
		t.Fatalf("expected comparison personal_best, got %s", settings.Comparison)
	}
}

func TestSettingsRoundTrip(t *testing.T) {
	store := tempStore(t)

	settings := DefaultSettings()
	settings.AlwaysOnTop = true
	settings.Hotkeys.StartSplit = "Enter"
	settings.Hotkeys.Pause = "KeyQ"

	if err := store.SaveSettings(settings); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := store.LoadSettings()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.AlwaysOnTop != true {
		t.Fatal("expected alwaysOnTop true")
	}

	if loaded.Hotkeys.StartSplit != "Enter" {
		t.Fatalf("expected startSplit Enter, got %s", loaded.Hotkeys.StartSplit)
	}

	if loaded.Hotkeys.Pause != "KeyQ" {
		t.Fatalf("expected pause KeyQ, got %s", loaded.Hotkeys.Pause)
	}

	// Unchanged fields should still have defaults.
	if loaded.Hotkeys.Reset != "KeyR" {
		t.Fatalf("expected reset KeyR, got %s", loaded.Hotkeys.Reset)
	}
}

func TestSettingsPartialJSON(t *testing.T) {
	store := tempStore(t)

	// Write a partial JSON file (missing hotkeys entirely).
	partial := []byte(`{"alwaysOnTop": true}`)
	if err := os.WriteFile(store.settingsPath(), partial, 0o640); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	loaded, err := store.LoadSettings()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.AlwaysOnTop != true {
		t.Fatal("expected alwaysOnTop true from partial JSON")
	}

	// Hotkeys should be defaults since they were missing from JSON.
	defaults := DefaultSettings()
	if loaded.Hotkeys.StartSplit != defaults.Hotkeys.StartSplit {
		t.Fatalf("expected default startSplit %q, got %q", defaults.Hotkeys.StartSplit, loaded.Hotkeys.StartSplit)
	}

	if loaded.Comparison != defaults.Comparison {
		t.Fatalf("expected default comparison %q, got %q", defaults.Comparison, loaded.Comparison)
	}
}

func TestSettingsAtomicWrite(t *testing.T) {
	store := tempStore(t)

	if err := store.SaveSettings(DefaultSettings()); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	entries, err := os.ReadDir(store.baseDir)
	if err != nil {
		t.Fatalf("readdir failed: %v", err)
	}

	for _, e := range entries {
		name := e.Name()
		if len(name) > 4 && name[len(name)-4:] == ".tmp" {
			t.Fatalf("tmp file left behind: %s", name)
		}
	}
}
