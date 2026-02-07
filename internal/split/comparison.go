package split

// Delta represents the time difference for a segment compared to a reference.
type Delta struct {
	SegmentIndex int   `json:"segmentIndex"`
	DeltaMS      int64 `json:"deltaMs"`    // Cumulative delta. Positive = behind, negative = ahead.
	IsBestEver   bool  `json:"isBestEver"` // True if this is the best segment time ever.
	IsAhead      bool  `json:"isAhead"`    // True if ahead of comparison.
	GainedTime   bool  `json:"gainedTime"` // True if segment time < comparison segment time.
	Skipped      bool  `json:"skipped"`    // True if the segment was skipped.
}

// ComputeDelta computes the delta for a single segment against a comparison split.
func ComputeDelta(currentSplitMS, comparisonSplitMS int64) int64 {
	return currentSplitMS - comparisonSplitMS
}

// ComparisonSplits returns the reference splits for the given comparison type.
func ComparisonSplits(att *Attempts, comparison string) []int64 {
	switch comparison {
	case "best_segments":
		return att.BestSegmentsCumulative()
	case "average_segments":
		return att.AverageSplits()
	case "latest_run":
		return att.LatestCompletedSplits()
	default:
		return att.PersonalBestSplits()
	}
}

// ComputeSplitDeltas computes deltas for all completed segments against the given comparison.
func ComputeSplitDeltas(att *Attempts, currentSplitsMS []int64, comparison string) []Delta {
	pbSplits := ComparisonSplits(att, comparison)
	deltas := make([]Delta, len(currentSplitsMS))

	for i, splitMS := range currentSplitsMS {
		d := Delta{SegmentIndex: i}

		if splitMS == 0 {
			d.Skipped = true
			deltas[i] = d

			continue
		}

		// Compute individual segment time.
		var segTimeMS int64
		if i == 0 {
			segTimeMS = splitMS
		} else {
			prevSplit := currentSplitsMS[i-1]
			if prevSplit == 0 {
				// Previous was skipped; can't compute accurate segment time.
				d.Skipped = true
				deltas[i] = d

				continue
			}

			segTimeMS = splitMS - prevSplit
		}

		// Check if this is the best segment ever.
		if i < len(att.Segments) && att.Segments[i].BestSegmentMS > 0 {
			d.IsBestEver = segTimeMS < att.Segments[i].BestSegmentMS
		}

		// Compute delta against comparison.
		if pbSplits != nil && i < len(pbSplits) && pbSplits[i] > 0 {
			d.DeltaMS = ComputeDelta(splitMS, pbSplits[i])
			d.IsAhead = d.DeltaMS < 0

			// Compute comparison segment time to determine if runner gained time.
			var compSegTimeMS int64
			if i == 0 {
				compSegTimeMS = pbSplits[i]
			} else if pbSplits[i-1] > 0 {
				compSegTimeMS = pbSplits[i] - pbSplits[i-1]
			}

			if compSegTimeMS > 0 {
				d.GainedTime = segTimeMS < compSegTimeMS
			}
		}

		deltas[i] = d
	}

	return deltas
}

// UpdatePersonalBest updates the PB cumulative times and best segment times.
func UpdatePersonalBest(att *Attempts, splitTimesMS []int64) {
	if len(splitTimesMS) != len(att.Segments) {
		return
	}

	// A run with the final segment skipped has no valid final time.
	finalTime := splitTimesMS[len(splitTimesMS)-1]
	if finalTime == 0 {
		// Still update best segment times for non-skipped segments.
		updateBestSegments(att, splitTimesMS)

		return
	}

	// Check against stored PB on segments, not attempts (avoids comparing against self).
	currentPBFinal := att.Segments[len(att.Segments)-1].PersonalBestMS
	isNewPB := currentPBFinal == 0 || finalTime < currentPBFinal

	// Update best segment times regardless of PB.
	updateBestSegments(att, splitTimesMS)

	// Update PB cumulative splits, preserving values for skipped segments.
	if isNewPB {
		for i, splitMS := range splitTimesMS {
			if i < len(att.Segments) && splitMS > 0 {
				att.Segments[i].PersonalBestMS = splitMS
			}
		}
	}
}

// RecalculatePersonalBest recomputes all PB and best segment data from scratch using history.
func RecalculatePersonalBest(att *Attempts) {
	// Clear existing PB data.
	for i := range att.Segments {
		att.Segments[i].PersonalBestMS = 0
		att.Segments[i].BestSegmentMS = 0
	}

	// Replay all completed attempts.
	for _, a := range att.History {
		if !a.Completed || len(a.SplitTimesMS) != len(att.Segments) {
			continue
		}

		UpdatePersonalBest(att, a.SplitTimesMS)
	}
}

func updateBestSegments(att *Attempts, splitTimesMS []int64) {
	for i, splitMS := range splitTimesMS {
		if splitMS == 0 || i >= len(att.Segments) {
			continue
		}

		var segTimeMS int64
		if i == 0 {
			segTimeMS = splitMS
		} else if splitTimesMS[i-1] > 0 {
			segTimeMS = splitMS - splitTimesMS[i-1]
		} else {
			continue // Can't compute if previous was skipped.
		}

		if att.Segments[i].BestSegmentMS == 0 || segTimeMS < att.Segments[i].BestSegmentMS {
			att.Segments[i].BestSegmentMS = segTimeMS
		}
	}
}
