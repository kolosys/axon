package axon

import (
	"sync/atomic"
	"time"
)

// ConnectionState represents the state of a WebSocket connection
type ConnectionState int32

const (
	// StateDisconnected indicates the connection is not established
	StateDisconnected ConnectionState = iota
	// StateConnecting indicates a connection attempt is in progress
	StateConnecting
	// StateConnected indicates the connection is established and ready
	StateConnected
	// StateReconnecting indicates a reconnection attempt is in progress
	StateReconnecting
	// StateClosing indicates the connection is being closed
	StateClosing
	// StateClosed indicates the connection has been permanently closed
	StateClosed
)

// String returns the string representation of the connection state
func (s ConnectionState) String() string {
	switch s {
	case StateDisconnected:
		return "disconnected"
	case StateConnecting:
		return "connecting"
	case StateConnected:
		return "connected"
	case StateReconnecting:
		return "reconnecting"
	case StateClosing:
		return "closing"
	case StateClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// IsActive returns true if the state represents an active or transitioning state
func (s ConnectionState) IsActive() bool {
	return s == StateConnecting || s == StateConnected || s == StateReconnecting
}

// CanReconnect returns true if reconnection is allowed from this state
func (s ConnectionState) CanReconnect() bool {
	return s == StateDisconnected || s == StateReconnecting
}

// StateChange represents a state transition event
type StateChange struct {
	From      ConnectionState
	To        ConnectionState
	Time      time.Time
	Err       error  // Error that caused the transition (if any)
	Attempt   int    // Reconnection attempt number (if applicable)
	SessionID string // Current session identifier
}

// StateHandler is a callback for state change events
type StateHandler func(StateChange)

// stateManager manages connection state transitions
type stateManager struct {
	state     atomic.Int32
	sessionID atomic.Value // string
	handlers  []StateHandler
}

// newStateManager creates a new state manager
func newStateManager() *stateManager {
	sm := &stateManager{}
	sm.state.Store(int32(StateDisconnected))
	sm.sessionID.Store("")
	return sm
}

// State returns the current connection state
func (sm *stateManager) State() ConnectionState {
	return ConnectionState(sm.state.Load())
}

// SetSessionID sets the current session identifier
func (sm *stateManager) SetSessionID(id string) {
	sm.sessionID.Store(id)
}

// SessionID returns the current session identifier
func (sm *stateManager) SessionID() string {
	v := sm.sessionID.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

// transition attempts to transition from one state to another
// Returns true if the transition was successful
func (sm *stateManager) transition(from, to ConnectionState, err error, attempt int) bool {
	if !sm.state.CompareAndSwap(int32(from), int32(to)) {
		return false
	}

	change := StateChange{
		From:      from,
		To:        to,
		Time:      time.Now(),
		Err:       err,
		Attempt:   attempt,
		SessionID: sm.SessionID(),
	}

	// Notify handlers
	for _, h := range sm.handlers {
		if h != nil {
			h(change)
		}
	}

	return true
}

// forceTransition transitions to a new state regardless of current state
func (sm *stateManager) forceTransition(to ConnectionState, err error, attempt int) ConnectionState {
	from := ConnectionState(sm.state.Swap(int32(to)))

	if from != to {
		change := StateChange{
			From:      from,
			To:        to,
			Time:      time.Now(),
			Err:       err,
			Attempt:   attempt,
			SessionID: sm.SessionID(),
		}

		// Notify handlers
		for _, h := range sm.handlers {
			if h != nil {
				h(change)
			}
		}
	}

	return from
}

// OnStateChange registers a callback for state change events
func (sm *stateManager) OnStateChange(handler StateHandler) {
	if handler != nil {
		sm.handlers = append(sm.handlers, handler)
	}
}

// validTransitions defines which state transitions are allowed
var validTransitions = map[ConnectionState][]ConnectionState{
	StateDisconnected: {StateConnecting, StateClosed},
	StateConnecting:   {StateConnected, StateDisconnected, StateClosed},
	StateConnected:    {StateDisconnected, StateReconnecting, StateClosing, StateClosed},
	StateReconnecting: {StateConnecting, StateConnected, StateDisconnected, StateClosed},
	StateClosing:      {StateClosed},
	StateClosed:       {}, // Terminal state
}

// isValidTransition checks if a state transition is valid
func isValidTransition(from, to ConnectionState) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
