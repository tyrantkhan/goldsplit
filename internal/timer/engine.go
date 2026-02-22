package timer

import (
	"sync"
	"time"
)

// TickData is emitted on each timer tick.
type TickData struct {
	ElapsedMS      int64    `json:"elapsedMs"`
	State          string   `json:"state"`
	CurrentSegment int      `json:"currentSegment"`
	SplitTimesMS   []int64  `json:"splitTimesMs"`
	SegmentTimesMS []int64  `json:"segmentTimesMs"`
	SplitNames     []string `json:"splitNames"`
}

// OnTickFunc is called on every timer tick with current data.
type OnTickFunc func(TickData)

// OnStateChangeFunc is called when the timer state changes.
type OnStateChangeFunc func(State)

// Engine is a high-precision speedrun timer.
type Engine struct {
	mu sync.RWMutex

	state      State
	startTime  time.Time
	pauseTime  time.Time
	pauseAccum time.Duration

	segmentNames   []string
	splitTimesMS   []int64 // cumulative split times in ms
	segmentTimesMS []int64 // individual segment durations in ms
	currentSegment int

	ticker   *time.Ticker
	stopChan chan struct{}

	onTick        OnTickFunc
	onStateChange OnStateChangeFunc
}

// New creates a new timer engine with the given segment names.
func New(segmentNames []string, onTick OnTickFunc, onStateChange OnStateChangeFunc) *Engine {
	return &Engine{
		segmentNames:  segmentNames,
		onTick:        onTick,
		onStateChange: onStateChange,
	}
}

// SetSegments replaces the segment list (only valid in Idle state).
func (e *Engine) SetSegments(names []string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state != Idle {
		return
	}

	e.segmentNames = names
}

// State returns the current timer state.
func (e *Engine) CurrentState() State {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.state
}

// ElapsedMS returns the current elapsed time in milliseconds.
func (e *Engine) ElapsedMS() int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.elapsedMS()
}

func (e *Engine) elapsedMS() int64 {
	switch e.state {
	case Idle:
		return 0
	case Running:
		return time.Since(e.startTime).Milliseconds() - e.pauseAccum.Milliseconds()
	case Paused:
		return e.pauseTime.Sub(e.startTime).Milliseconds() - e.pauseAccum.Milliseconds()
	case Finished:
		if len(e.splitTimesMS) > 0 {
			return e.splitTimesMS[len(e.splitTimesMS)-1]
		}

		return 0
	}

	return 0
}

// Start begins the timer. Only valid from Idle state.
func (e *Engine) Start() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state != Idle {
		return
	}

	e.state = Running
	e.startTime = time.Now()
	e.pauseAccum = 0
	e.currentSegment = 0
	e.splitTimesMS = nil
	e.segmentTimesMS = nil

	e.startTicker()
	e.notifyStateChange()
}

// Split records the current segment time. Only valid from Running state.
func (e *Engine) Split() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state != Running {
		return
	}

	elapsed := e.elapsedMS()
	e.splitTimesMS = append(e.splitTimesMS, elapsed)

	var segTime int64
	if len(e.splitTimesMS) == 1 {
		segTime = elapsed
	} else {
		segTime = elapsed - e.splitTimesMS[len(e.splitTimesMS)-2]
	}

	e.segmentTimesMS = append(e.segmentTimesMS, segTime)
	e.currentSegment++

	if e.currentSegment >= len(e.segmentNames) {
		e.state = Finished
		e.stopTicker()
		e.notifyTick()
		e.notifyStateChange()
	}
}

// SkipSplit skips the current segment without recording a time.
func (e *Engine) SkipSplit() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state != Running {
		return
	}

	e.splitTimesMS = append(e.splitTimesMS, 0) // 0 indicates skipped
	e.segmentTimesMS = append(e.segmentTimesMS, 0)
	e.currentSegment++

	if e.currentSegment >= len(e.segmentNames) {
		e.state = Finished
		e.stopTicker()
		e.notifyTick()
		e.notifyStateChange()
	}
}

// UndoSplit reverts the last split. Only valid from Running state.
func (e *Engine) UndoSplit() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state != Running || e.currentSegment == 0 {
		return
	}

	e.splitTimesMS = e.splitTimesMS[:len(e.splitTimesMS)-1]
	e.segmentTimesMS = e.segmentTimesMS[:len(e.segmentTimesMS)-1]
	e.currentSegment--
	e.notifyTick()
}

// Pause pauses the timer. Only valid from Running state.
func (e *Engine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state != Running {
		return
	}

	e.state = Paused
	e.pauseTime = time.Now()
	e.stopTicker()
	e.notifyStateChange()
}

// Resume resumes the timer from pause. Only valid from Paused state.
func (e *Engine) Resume() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state != Paused {
		return
	}

	e.pauseAccum += time.Since(e.pauseTime)
	e.state = Running
	e.startTicker()
	e.notifyStateChange()
}

// Restore puts the engine into Paused state with previously saved data.
// Only valid from Idle state. After restoring, the normal Resume transition works.
func (e *Engine) Restore(elapsedMS int64, currentSegment int, splitTimesMS, segmentTimesMS []int64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state != Idle {
		return
	}

	now := time.Now()
	e.state = Paused
	e.startTime = now.Add(-time.Duration(elapsedMS) * time.Millisecond)
	e.pauseTime = now
	e.pauseAccum = 0
	e.currentSegment = currentSegment
	e.splitTimesMS = make([]int64, len(splitTimesMS))
	copy(e.splitTimesMS, splitTimesMS)
	e.segmentTimesMS = make([]int64, len(segmentTimesMS))
	copy(e.segmentTimesMS, segmentTimesMS)

	e.notifyTick()
	e.notifyStateChange()
}

// Reset stops the timer and returns to Idle. Valid from any state except Idle.
func (e *Engine) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state == Idle {
		return
	}

	e.stopTicker()
	e.state = Idle
	e.currentSegment = 0
	e.splitTimesMS = nil
	e.segmentTimesMS = nil
	e.notifyTick()
	e.notifyStateChange()
}

// GetTickData returns the current tick data snapshot.
func (e *Engine) GetTickData() TickData {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.tickData()
}

func (e *Engine) tickData() TickData {
	splitsCopy := make([]int64, len(e.splitTimesMS))
	copy(splitsCopy, e.splitTimesMS)

	segsCopy := make([]int64, len(e.segmentTimesMS))
	copy(segsCopy, e.segmentTimesMS)

	namesCopy := make([]string, len(e.segmentNames))
	copy(namesCopy, e.segmentNames)

	return TickData{
		ElapsedMS:      e.elapsedMS(),
		State:          e.state.String(),
		CurrentSegment: e.currentSegment,
		SplitTimesMS:   splitsCopy,
		SegmentTimesMS: segsCopy,
		SplitNames:     namesCopy,
	}
}

func (e *Engine) startTicker() {
	e.stopChan = make(chan struct{})
	e.ticker = time.NewTicker(15 * time.Millisecond)

	go func() {
		for {
			select {
			case <-e.ticker.C:
				e.mu.RLock()
				data := e.tickData()
				e.mu.RUnlock()

				if e.onTick != nil {
					e.onTick(data)
				}
			case <-e.stopChan:
				return
			}
		}
	}()
}

func (e *Engine) stopTicker() {
	if e.ticker != nil {
		e.ticker.Stop()
	}

	if e.stopChan != nil {
		close(e.stopChan)
		e.stopChan = nil
	}
}

func (e *Engine) notifyTick() {
	if e.onTick != nil {
		e.onTick(e.tickData())
	}
}

func (e *Engine) notifyStateChange() {
	if e.onStateChange != nil {
		e.onStateChange(e.state)
	}
}

// SplitTimesMS returns a copy of the recorded split times.
func (e *Engine) SplitTimesMS() []int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]int64, len(e.splitTimesMS))
	copy(result, e.splitTimesMS)

	return result
}

// SegmentTimesMS returns a copy of the recorded segment times.
func (e *Engine) SegmentTimesMS() []int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]int64, len(e.segmentTimesMS))
	copy(result, e.segmentTimesMS)

	return result
}

// CurrentSegment returns the current segment index.
func (e *Engine) CurrentSegment() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.currentSegment
}
