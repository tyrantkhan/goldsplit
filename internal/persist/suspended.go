package persist

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// SuspendedRun holds the state of an in-progress run that was suspended.
type SuspendedRun struct {
	TemplateID     string  `json:"templateId"`
	AttemptsID     string  `json:"attemptsId"`
	ElapsedMS      int64   `json:"elapsedMs"`
	CurrentSegment int     `json:"currentSegment"`
	SplitTimesMS   []int64 `json:"splitTimesMs"`
	SegmentTimesMS []int64 `json:"segmentTimesMs"`
	SuspendedAt    int64   `json:"suspendedAt"`
}

// SaveSuspendedRun persists a suspended run to disk using atomic write.
func (s *Store) SaveSuspendedRun(run *SuspendedRun) error {
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling suspended run: %w", err)
	}

	path := s.suspendedRunPath()
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

// LoadSuspendedRun reads a suspended run from disk.
// Returns (nil, nil) if the file does not exist.
func (s *Store) LoadSuspendedRun() (*SuspendedRun, error) {
	data, err := os.ReadFile(s.suspendedRunPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("reading suspended run file: %w", err)
	}

	var run SuspendedRun
	if err := json.Unmarshal(data, &run); err != nil {
		return nil, fmt.Errorf("unmarshaling suspended run: %w", err)
	}

	return &run, nil
}

// DeleteSuspendedRun removes the suspended run file. No-op if the file does not exist.
func (s *Store) DeleteSuspendedRun() error {
	err := os.Remove(s.suspendedRunPath())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("deleting suspended run file: %w", err)
	}

	return nil
}

func (s *Store) suspendedRunPath() string {
	return filepath.Join(s.baseDir, "suspended_run.json")
}
