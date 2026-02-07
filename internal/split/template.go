package split

import "time"

// Template is a reusable blueprint for a speedrun: game name and segment names.
type Template struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	SegmentNames []string  `json:"segmentNames"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// NewTemplate creates a new template with the given parameters.
func NewTemplate(id, name string, segmentNames []string) *Template {
	now := time.Now()

	return &Template{
		ID:           id,
		Name:         name,
		SegmentNames: segmentNames,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
