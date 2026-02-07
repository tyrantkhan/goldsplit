package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"goldsplit/internal/split"
)

// Store handles persistence of templates and attempts to disk.
type Store struct {
	baseDir string
}

// NewStore creates a store at the given base directory.
func NewStore(baseDir string) (*Store, error) {
	for _, sub := range []string{"templates", "attempts"} {
		dir := filepath.Join(baseDir, sub)
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return nil, fmt.Errorf("creating %s directory: %w", sub, err)
		}
	}

	return &Store{baseDir: baseDir}, nil
}

// DefaultBaseDir returns the default application data directory.
func DefaultBaseDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	return filepath.Join(home, "Library", "Application Support", "goldsplit"), nil
}

// TemplateSummary is a lightweight representation of a template for listing.
type TemplateSummary struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	SegmentCount int    `json:"segmentCount"`
	UpdatedAt    int64  `json:"updatedAt"`
}

// AttemptsSummary is a lightweight representation of attempts for listing.
type AttemptsSummary struct {
	ID           string `json:"id"`
	TemplateID   string `json:"templateId"`
	Name         string `json:"name"`
	CategoryName string `json:"categoryName"`
	AttemptCount int    `json:"attemptCount"`
	UpdatedAt    int64  `json:"updatedAt"`
}

// SaveTemplate persists a template to disk using atomic write.
func (s *Store) SaveTemplate(tmpl *split.Template) error {
	data, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling template: %w", err)
	}

	path := s.templatePath(tmpl.ID)
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

// LoadTemplate reads a template from disk by ID.
func (s *Store) LoadTemplate(id string) (*split.Template, error) {
	data, err := os.ReadFile(s.templatePath(id))
	if err != nil {
		return nil, fmt.Errorf("reading template file: %w", err)
	}

	var tmpl split.Template
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("unmarshaling template: %w", err)
	}

	return &tmpl, nil
}

// ListTemplates returns all saved templates.
func (s *Store) ListTemplates() ([]TemplateSummary, error) {
	dir := filepath.Join(s.baseDir, "templates")

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading templates directory: %w", err)
	}

	var summaries []TemplateSummary

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".json")

		tmpl, err := s.LoadTemplate(id)
		if err != nil {
			continue
		}

		summaries = append(summaries, TemplateSummary{
			ID:           tmpl.ID,
			Name:         tmpl.Name,
			SegmentCount: len(tmpl.SegmentNames),
			UpdatedAt:    tmpl.UpdatedAt.Unix(),
		})
	}

	return summaries, nil
}

// DeleteTemplate removes a template and all associated attempts from disk.
func (s *Store) DeleteTemplate(id string) error {
	// Delete associated attempts first.
	attDir := filepath.Join(s.baseDir, "attempts")

	entries, err := os.ReadDir(attDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}

			attID := strings.TrimSuffix(entry.Name(), ".json")

			att, loadErr := s.LoadAttempts(attID)
			if loadErr != nil {
				continue
			}

			if att.TemplateID == id {
				_ = os.Remove(s.attemptsPath(attID))
			}
		}
	}

	if err := os.Remove(s.templatePath(id)); err != nil {
		return fmt.Errorf("deleting template file: %w", err)
	}

	return nil
}

// SaveAttempts persists attempts to disk using atomic write.
func (s *Store) SaveAttempts(att *split.Attempts) error {
	data, err := json.MarshalIndent(att, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling attempts: %w", err)
	}

	path := s.attemptsPath(att.ID)
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

// LoadAttempts reads attempts from disk by ID.
func (s *Store) LoadAttempts(id string) (*split.Attempts, error) {
	data, err := os.ReadFile(s.attemptsPath(id))
	if err != nil {
		return nil, fmt.Errorf("reading attempts file: %w", err)
	}

	var att split.Attempts
	if err := json.Unmarshal(data, &att); err != nil {
		return nil, fmt.Errorf("unmarshaling attempts: %w", err)
	}

	return &att, nil
}

// ListAttemptsForTemplate returns all attempts for a given template.
func (s *Store) ListAttemptsForTemplate(templateID string) ([]AttemptsSummary, error) {
	dir := filepath.Join(s.baseDir, "attempts")

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading attempts directory: %w", err)
	}

	var summaries []AttemptsSummary

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".json")

		att, err := s.LoadAttempts(id)
		if err != nil {
			continue
		}

		if att.TemplateID != templateID {
			continue
		}

		summaries = append(summaries, AttemptsSummary{
			ID:           att.ID,
			TemplateID:   att.TemplateID,
			Name:         att.Name,
			CategoryName: att.CategoryName,
			AttemptCount: att.AttemptCount,
			UpdatedAt:    att.UpdatedAt.Unix(),
		})
	}

	return summaries, nil
}

// DeleteAttempts removes an attempts file from disk.
func (s *Store) DeleteAttempts(id string) error {
	if err := os.Remove(s.attemptsPath(id)); err != nil {
		return fmt.Errorf("deleting attempts file: %w", err)
	}

	return nil
}

func (s *Store) templatePath(id string) string {
	return filepath.Join(s.baseDir, "templates", id+".json")
}

func (s *Store) attemptsPath(id string) string {
	return filepath.Join(s.baseDir, "attempts", id+".json")
}
