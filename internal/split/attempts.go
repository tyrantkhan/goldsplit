package split

import "time"

// Segment represents a single segment in a run.
type Segment struct {
	Name           string `json:"name"`
	PersonalBestMS int64  `json:"personalBestMs"` // Cumulative PB time.
	BestSegmentMS  int64  `json:"bestSegmentMs"`  // Best individual segment time ever.
}

// Attempt records a single attempt.
type Attempt struct {
	ID           int       `json:"id"`
	StartedAt    time.Time `json:"startedAt"`
	SplitTimesMS []int64   `json:"splitTimesMs"` // Cumulative split times (0 = skipped).
	Completed    bool      `json:"completed"`
}

// Attempts tracks category-specific data: segments (snapshotted from a template),
// PB/best segment data, and attempt history.
type Attempts struct {
	ID           string    `json:"id"`
	TemplateID   string    `json:"templateId"`
	Name         string    `json:"name"`
	CategoryName string    `json:"categoryName"`
	Segments     []Segment `json:"segments"`
	AttemptCount int       `json:"attemptCount"`
	History      []Attempt `json:"history"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// NewAttempts creates a new Attempts with segments snapshotted from segment names.
func NewAttempts(id, templateID, name, categoryName string, segmentNames []string) *Attempts {
	segs := make([]Segment, len(segmentNames))
	for i, n := range segmentNames {
		segs[i] = Segment{Name: n}
	}

	now := time.Now()

	return &Attempts{
		ID:           id,
		TemplateID:   templateID,
		Name:         name,
		CategoryName: categoryName,
		Segments:     segs,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// SegmentNames returns the names of all segments.
func (a *Attempts) SegmentNames() []string {
	names := make([]string, len(a.Segments))
	for i, s := range a.Segments {
		names[i] = s.Name
	}

	return names
}

// AddAttempt records a new attempt.
func (a *Attempts) AddAttempt(splitTimesMS []int64, completed bool) {
	a.AttemptCount++
	a.History = append(a.History, Attempt{
		ID:           a.AttemptCount,
		StartedAt:    time.Now(),
		SplitTimesMS: splitTimesMS,
		Completed:    completed,
	})
	a.UpdatedAt = time.Now()
}

// PersonalBestSplits returns the split times from the PB attempt, or nil if no PB exists.
// Attempts with skipped final segments (finalTime == 0) are excluded.
func (a *Attempts) PersonalBestSplits() []int64 {
	var bestTime int64
	var bestSplits []int64

	for _, att := range a.History {
		if !att.Completed || len(att.SplitTimesMS) == 0 {
			continue
		}

		finalTime := att.SplitTimesMS[len(att.SplitTimesMS)-1]
		if finalTime == 0 {
			continue
		}

		if bestSplits == nil || finalTime < bestTime {
			bestTime = finalTime
			bestSplits = make([]int64, len(att.SplitTimesMS))
			copy(bestSplits, att.SplitTimesMS)
		}
	}

	return bestSplits
}

// BestSegmentsCumulative returns cumulative splits built from each segment's best time.
// Returns nil if any segment has no best time recorded.
func (a *Attempts) BestSegmentsCumulative() []int64 {
	cumulative := make([]int64, len(a.Segments))
	var sum int64

	for i, seg := range a.Segments {
		if seg.BestSegmentMS == 0 {
			return nil
		}

		sum += seg.BestSegmentMS
		cumulative[i] = sum
	}

	return cumulative
}

// AverageSplits returns the average cumulative split time at each segment across completed attempts.
// Skips attempts with finalTime==0 and per-segment values of 0. Returns nil if no valid data.
func (a *Attempts) AverageSplits() []int64 {
	n := len(a.Segments)
	sums := make([]int64, n)
	counts := make([]int, n)

	for _, att := range a.History {
		if !att.Completed || len(att.SplitTimesMS) != n {
			continue
		}

		finalTime := att.SplitTimesMS[n-1]
		if finalTime == 0 {
			continue
		}

		for i, t := range att.SplitTimesMS {
			if t == 0 {
				continue
			}

			sums[i] += t
			counts[i]++
		}
	}

	result := make([]int64, n)
	hasData := false

	for i := range result {
		if counts[i] == 0 {
			result[i] = 0

			continue
		}

		result[i] = sums[i] / int64(counts[i])
		hasData = true
	}

	if !hasData {
		return nil
	}

	return result
}

// LatestCompletedSplits returns the splits from the most recent completed attempt.
// Excludes attempts with finalTime==0. Returns nil if none.
func (a *Attempts) LatestCompletedSplits() []int64 {
	for i := len(a.History) - 1; i >= 0; i-- {
		att := a.History[i]
		if !att.Completed || len(att.SplitTimesMS) == 0 {
			continue
		}

		finalTime := att.SplitTimesMS[len(att.SplitTimesMS)-1]
		if finalTime == 0 {
			continue
		}

		result := make([]int64, len(att.SplitTimesMS))
		copy(result, att.SplitTimesMS)

		return result
	}

	return nil
}

// DeleteAttempt removes an attempt by ID from History and recalculates PB data.
func (a *Attempts) DeleteAttempt(attemptID int) bool {
	idx := -1
	for i, att := range a.History {
		if att.ID == attemptID {
			idx = i

			break
		}
	}

	if idx == -1 {
		return false
	}

	a.History = append(a.History[:idx], a.History[idx+1:]...)
	a.AttemptCount--
	a.UpdatedAt = time.Now()
	RecalculatePersonalBest(a)

	return true
}

// EditAttemptSplits updates split times for an attempt and recalculates PB data.
// Returns false if the attempt is not found, the segment count mismatches,
// or cumulative times are not monotonically increasing (for non-zero values).
func (a *Attempts) EditAttemptSplits(attemptID int, newSplits []int64) bool {
	if len(newSplits) != len(a.Segments) {
		return false
	}

	// Validate monotonically increasing for non-zero values.
	var lastNonZero int64

	for _, v := range newSplits {
		if v == 0 {
			continue
		}

		if v <= lastNonZero {
			return false
		}

		lastNonZero = v
	}

	for i := range a.History {
		if a.History[i].ID == attemptID {
			a.History[i].SplitTimesMS = newSplits
			a.UpdatedAt = time.Now()
			RecalculatePersonalBest(a)

			return true
		}
	}

	return false
}
