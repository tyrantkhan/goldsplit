package hotkey

// Action represents a hotkey action.
type Action int

const (
	ActionStartSplit Action = iota // Context-dependent: start if idle, split if running.
	ActionPause                    // Pause/unpause.
	ActionReset                    // Reset the timer.
	ActionUndoSplit                // Undo the last split.
	ActionSkipSplit                // Skip the current segment.
)

func (a Action) String() string {
	switch a {
	case ActionStartSplit:
		return "start_split"
	case ActionPause:
		return "pause"
	case ActionReset:
		return "reset"
	case ActionUndoSplit:
		return "undo_split"
	case ActionSkipSplit:
		return "skip_split"
	default:
		return "unknown"
	}
}

// Handler is called when a hotkey action fires.
type Handler func(Action)

// Manager manages global hotkey registration.
// Note: Global hotkey registration via golang.design/x/hotkey requires CGO
// and platform-specific setup. For the MVP, hotkeys are handled via the
// frontend (keyboard events) and the Wails-bound methods. True global
// hotkeys (when window is unfocused) can be added as an enhancement.
type Manager struct {
	handler Handler
	enabled bool
}

// NewManager creates a hotkey manager with the given handler.
func NewManager(handler Handler) *Manager {
	return &Manager{
		handler: handler,
		enabled: true,
	}
}

// Start begins listening for hotkeys. In the MVP, this is a no-op since
// hotkeys are handled via frontend keyboard events calling Go methods.
func (m *Manager) Start() error {
	return nil
}

// Stop stops listening for hotkeys.
func (m *Manager) Stop() {
	m.enabled = false
}

// Dispatch manually triggers a hotkey action (used by frontend bridge).
func (m *Manager) Dispatch(action Action) {
	if m.enabled && m.handler != nil {
		m.handler(action)
	}
}
