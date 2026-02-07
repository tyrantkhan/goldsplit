package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ColorSettings holds the hex color values for delta display.
type ColorSettings struct {
	AheadGaining  string `json:"aheadGaining"`
	AheadLosing   string `json:"aheadLosing"`
	BehindGaining string `json:"behindGaining"`
	BehindLosing  string `json:"behindLosing"`
	BestSegment   string `json:"bestSegment"`
}

// Settings holds the application settings.
type Settings struct {
	AlwaysOnTop bool           `json:"alwaysOnTop"`
	Hotkeys     HotkeyBindings `json:"hotkeys"`
	Comparison  string         `json:"comparison"`
	Colors      ColorSettings  `json:"colors"`
}

// HotkeyBindings holds the key bindings for each action.
// Values are JS KeyboardEvent.code strings (e.g. "Space", "KeyP").
type HotkeyBindings struct {
	StartSplit string `json:"startSplit"`
	Pause      string `json:"pause"`
	Reset      string `json:"reset"`
	UndoSplit  string `json:"undoSplit"`
	SkipSplit  string `json:"skipSplit"`
}

// DefaultSettings returns the default settings for a fresh install.
func DefaultSettings() Settings {
	return Settings{
		AlwaysOnTop: false,
		Hotkeys: HotkeyBindings{
			StartSplit: "Space",
			Pause:      "KeyP",
			Reset:      "KeyR",
			UndoSplit:  "Backspace",
			SkipSplit:  "KeyS",
		},
		Comparison: "personal_best",
		Colors: ColorSettings{
			AheadGaining:  "#30d158",
			AheadLosing:   "#7ec890",
			BehindGaining: "#cc6b65",
			BehindLosing:  "#ff453a",
			BestSegment:   "#ffd60a",
		},
	}
}

// LoadSettings reads settings from disk. Returns defaults if the file is missing.
func (s *Store) LoadSettings() (Settings, error) {
	settings := DefaultSettings()

	data, err := os.ReadFile(s.settingsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return settings, nil
		}

		return settings, fmt.Errorf("reading settings file: %w", err)
	}

	// Unmarshal into pre-populated defaults so new fields auto-fill.
	if err := json.Unmarshal(data, &settings); err != nil {
		return DefaultSettings(), fmt.Errorf("unmarshaling settings: %w", err)
	}

	return settings, nil
}

// SaveSettings persists settings to disk using atomic write.
func (s *Store) SaveSettings(settings Settings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling settings: %w", err)
	}

	path := s.settingsPath()
	tmpPath := path + ".tmp"

	if err := os.WriteFile(tmpPath, data, 0o640); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)

		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}

func (s *Store) settingsPath() string {
	return filepath.Join(s.baseDir, "settings.json")
}
