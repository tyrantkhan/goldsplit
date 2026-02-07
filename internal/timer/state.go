package timer

// State represents the current state of the timer.
type State int

const (
	Idle     State = iota // Timer has not started or has been reset.
	Running               // Timer is actively counting.
	Paused                // Timer is paused mid-run.
	Finished              // Run completed (all segments split).
)

func (s State) String() string {
	switch s {
	case Idle:
		return "idle"
	case Running:
		return "running"
	case Paused:
		return "paused"
	case Finished:
		return "finished"
	default:
		return "unknown"
	}
}
