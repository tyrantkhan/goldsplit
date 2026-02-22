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
		return att.LatestRunSplits()
	default:
		return att.PersonalBestSplits()
	}
}

// ComputeSplitDeltas computes deltas for all completed segments against the given comparison.
func ComputeSplitDeltas(att *Attempts, currentSplitsMS []int64, comparison string) []Delta {
	compSplits := ComparisonSplits(att, comparison)
	bestSegs := att.BestSegments()
	deltas := make([]Delta, len(currentSplitsMS))

	for i, splitMS := range currentSplitsMS {
		d := Delta{SegmentIndex: i}

		if splitMS == 0 {
			d.Skipped = true
			deltas[i] = d

			continue
		}

		// Compute individual segment time by finding last non-skipped cumulative.
		var segTimeMS int64
		if i == 0 {
			segTimeMS = splitMS
		} else {
			prevSplit, found := lastNonZeroBefore(currentSplitsMS, i)
			if !found {
				d.Skipped = true
				deltas[i] = d

				continue
			}

			segTimeMS = splitMS - prevSplit
		}

		// Check if this is the best segment ever.
		if i < len(bestSegs) && bestSegs[i] > 0 {
			d.IsBestEver = segTimeMS < bestSegs[i]
		}

		// Compute delta against comparison.
		if compSplits != nil && i < len(compSplits) && compSplits[i] > 0 {
			d.DeltaMS = ComputeDelta(splitMS, compSplits[i])
			d.IsAhead = d.DeltaMS < 0

			// Compute comparison segment time to determine if runner gained time.
			var compSegTimeMS int64
			if i == 0 {
				compSegTimeMS = compSplits[i]
			} else if prevComp, ok := lastNonZeroBefore(compSplits, i); ok {
				compSegTimeMS = compSplits[i] - prevComp
			}

			if compSegTimeMS > 0 {
				d.GainedTime = segTimeMS < compSegTimeMS
			}
		}

		deltas[i] = d
	}

	return deltas
}

// lastNonZeroBefore returns the most recent non-zero cumulative split before index i.
func lastNonZeroBefore(splits []int64, i int) (int64, bool) {
	for j := i - 1; j >= 0; j-- {
		if splits[j] != 0 {
			return splits[j], true
		}
	}

	return 0, false
}
