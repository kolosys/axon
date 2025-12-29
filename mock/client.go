package mock

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kolosys/axon"
)

// errorHolder wraps an error for use with atomic.Value
type errorHolder struct {
	err error
}

// MockClient provides an in-memory mock of axon.Client[T] for testing
type MockClient[T any] struct {
	state     atomic.Int32
	sessionID atomic.Value

	readCh  chan T
	writeCh chan T

	onConnect     func()
	onDisconnect  func(error)
	onMessage     func(T)
	onStateChange func(axon.StateChange)

	mu           sync.RWMutex
	writtenMsgs  []T
	receivedMsgs []T
	stateChanges []axon.StateChange

	connectErr atomic.Value
	readErr    atomic.Value
	writeErr   atomic.Value

	connectDelay time.Duration
	readDelay    time.Duration
	writeDelay   time.Duration

	queueStats atomic.Value

	closed    atomic.Bool
	closeOnce sync.Once
}

// MockClientOptions configures MockClient behavior
type MockClientOptions struct {
	BufferSize         int
	InitialState       axon.ConnectionState
	ConnectDelay       time.Duration
	ReadDelay          time.Duration
	WriteDelay         time.Duration
	RecordMessages     bool
	RecordStateChanges bool
	QueueSize          int
}

// NewMockClient creates a new mock client
func NewMockClient[T any](opts *MockClientOptions) *MockClient[T] {
	if opts == nil {
		opts = &MockClientOptions{
			BufferSize:         100,
			InitialState:       axon.StateDisconnected,
			RecordMessages:     true,
			RecordStateChanges: true,
			QueueSize:          100,
		}
	}

	if opts.BufferSize <= 0 {
		opts.BufferSize = 100
	}
	if opts.QueueSize <= 0 {
		opts.QueueSize = 100
	}

	mc := &MockClient[T]{
		readCh:       make(chan T, opts.BufferSize),
		writeCh:      make(chan T, opts.BufferSize),
		connectDelay: opts.ConnectDelay,
		readDelay:    opts.ReadDelay,
		writeDelay:   opts.WriteDelay,
	}

	mc.state.Store(int32(opts.InitialState))
	mc.sessionID.Store("")

	mc.queueStats.Store(axon.MessageQueueStats{
		MaxSize: opts.QueueSize,
	})

	if opts.RecordMessages {
		mc.writtenMsgs = make([]T, 0, 50)
		mc.receivedMsgs = make([]T, 0, 50)
	}
	if opts.RecordStateChanges {
		mc.stateChanges = make([]axon.StateChange, 0, 20)
	}

	return mc
}

// Connect simulates establishing a connection
func (mc *MockClient[T]) Connect(ctx context.Context) error {
	if eh := mc.connectErr.Load(); eh != nil {
		if holder := eh.(*errorHolder); holder != nil && holder.err != nil {
			return holder.err
		}
	}

	if !mc.transitionState(axon.StateDisconnected, axon.StateConnecting, nil, 0) {
		if mc.State() == axon.StateConnected {
			return nil
		}
		return axon.ErrInvalidState
	}

	if mc.connectDelay > 0 {
		select {
		case <-ctx.Done():
			mc.transitionState(axon.StateConnecting, axon.StateDisconnected, ctx.Err(), 0)
			return axon.ErrContextCanceled
		case <-time.After(mc.connectDelay):
		}
	}

	mc.transitionState(axon.StateConnecting, axon.StateConnected, nil, 0)

	if mc.onConnect != nil {
		mc.onConnect()
	}

	return nil
}

// ConnectWithReadLoop connects and starts message delivery
func (mc *MockClient[T]) ConnectWithReadLoop(ctx context.Context) error {
	if err := mc.Connect(ctx); err != nil {
		return err
	}

	go mc.readLoop(ctx)
	return nil
}

// readLoop continuously delivers messages via OnMessage callback
func (mc *MockClient[T]) readLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-mc.readCh:
			if !ok {
				return
			}

			if mc.onMessage != nil {
				mc.onMessage(msg)
			}
		}
	}
}

// Read reads a message from the client
func (mc *MockClient[T]) Read(ctx context.Context) (T, error) {
	var zero T

	if mc.State() != axon.StateConnected {
		return zero, axon.ErrConnectionClosed
	}

	if eh := mc.readErr.Load(); eh != nil {
		if holder := eh.(*errorHolder); holder != nil && holder.err != nil {
			return zero, holder.err
		}
	}

	if mc.readDelay > 0 {
		select {
		case <-ctx.Done():
			return zero, axon.ErrContextCanceled
		case <-time.After(mc.readDelay):
		}
	}

	select {
	case <-ctx.Done():
		return zero, axon.ErrContextCanceled
	case msg, ok := <-mc.readCh:
		if !ok {
			return zero, axon.ErrConnectionClosed
		}

		if mc.receivedMsgs != nil {
			mc.mu.Lock()
			mc.receivedMsgs = append(mc.receivedMsgs, msg)
			mc.mu.Unlock()
		}

		return msg, nil
	}
}

// Write writes a message to the client
func (mc *MockClient[T]) Write(ctx context.Context, msg T) error {
	state := mc.State()

	if state != axon.StateConnected {
		return axon.ErrConnectionClosed
	}

	if eh := mc.writeErr.Load(); eh != nil {
		if holder := eh.(*errorHolder); holder != nil && holder.err != nil {
			return holder.err
		}
	}

	if mc.writeDelay > 0 {
		select {
		case <-ctx.Done():
			return axon.ErrContextCanceled
		case <-time.After(mc.writeDelay):
		}
	}

	if mc.writtenMsgs != nil {
		mc.mu.Lock()
		mc.writtenMsgs = append(mc.writtenMsgs, msg)
		mc.mu.Unlock()
	}

	select {
	case <-ctx.Done():
		return axon.ErrContextCanceled
	case mc.writeCh <- msg:
		return nil
	}
}

// State returns the current connection state
func (mc *MockClient[T]) State() axon.ConnectionState {
	return axon.ConnectionState(mc.state.Load())
}

// IsConnected returns true if the client is connected
func (mc *MockClient[T]) IsConnected() bool {
	return mc.State() == axon.StateConnected
}

// Close closes the mock client
func (mc *MockClient[T]) Close() error {
	mc.closeOnce.Do(func() {
		old := mc.State()
		if old != axon.StateClosed {
			mc.transitionState(old, axon.StateClosing, nil, 0)
		}
		mc.closed.Store(true)

		close(mc.readCh)
		close(mc.writeCh)

		if old != axon.StateClosed {
			mc.transitionState(axon.StateClosing, axon.StateClosed, nil, 0)
		}
	})

	return nil
}

// OnConnect sets the connect callback
func (mc *MockClient[T]) OnConnect(fn func()) {
	mc.onConnect = fn
}

// OnDisconnect sets the disconnect callback
func (mc *MockClient[T]) OnDisconnect(fn func(error)) {
	mc.onDisconnect = fn
}

// OnMessage sets the message callback
func (mc *MockClient[T]) OnMessage(fn func(T)) {
	mc.onMessage = fn
}

// OnStateChange sets the state change callback
func (mc *MockClient[T]) OnStateChange(fn func(axon.StateChange)) {
	mc.onStateChange = fn
}

// SetSessionID sets the session ID
func (mc *MockClient[T]) SetSessionID(id string) {
	mc.sessionID.Store(id)
}

// SessionID returns the session ID
func (mc *MockClient[T]) SessionID() string {
	v := mc.sessionID.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

// QueueStats returns queue statistics
func (mc *MockClient[T]) QueueStats() axon.MessageQueueStats {
	v := mc.queueStats.Load()
	if v == nil {
		return axon.MessageQueueStats{}
	}
	return v.(axon.MessageQueueStats)
}

// transitionState performs a state transition
func (mc *MockClient[T]) transitionState(from, to axon.ConnectionState, err error, attempt int) bool {
	if !mc.state.CompareAndSwap(int32(from), int32(to)) {
		return false
	}

	change := axon.StateChange{
		From:      from,
		To:        to,
		Time:      time.Now(),
		Err:       err,
		Attempt:   attempt,
		SessionID: mc.SessionID(),
	}

	if mc.stateChanges != nil {
		mc.mu.Lock()
		mc.stateChanges = append(mc.stateChanges, change)
		mc.mu.Unlock()
	}

	if mc.onStateChange != nil {
		mc.onStateChange(change)
	}

	return true
}

// Test Helpers

// InjectMessage simulates receiving a message
func (mc *MockClient[T]) InjectMessage(msg T) error {
	if mc.closed.Load() {
		return axon.ErrConnectionClosed
	}

	select {
	case mc.readCh <- msg:
		return nil
	default:
		return axon.ErrQueueFull
	}
}

// DrainWritten retrieves all written messages
func (mc *MockClient[T]) DrainWritten() []T {
	var msgs []T
	for {
		select {
		case msg := <-mc.writeCh:
			msgs = append(msgs, msg)
		default:
			return msgs
		}
	}
}

// WrittenMessages returns recorded written messages
func (mc *MockClient[T]) WrittenMessages() []T {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.writtenMsgs == nil {
		return nil
	}

	msgs := make([]T, len(mc.writtenMsgs))
	copy(msgs, mc.writtenMsgs)
	return msgs
}

// ReceivedMessages returns recorded received messages
func (mc *MockClient[T]) ReceivedMessages() []T {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.receivedMsgs == nil {
		return nil
	}

	msgs := make([]T, len(mc.receivedMsgs))
	copy(msgs, mc.receivedMsgs)
	return msgs
}

// StateChanges returns recorded state changes
func (mc *MockClient[T]) StateChanges() []axon.StateChange {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.stateChanges == nil {
		return nil
	}

	changes := make([]axon.StateChange, len(mc.stateChanges))
	copy(changes, mc.stateChanges)
	return changes
}

// ClearRecorded clears all recorded data
func (mc *MockClient[T]) ClearRecorded() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.writtenMsgs != nil {
		mc.writtenMsgs = mc.writtenMsgs[:0]
	}
	if mc.receivedMsgs != nil {
		mc.receivedMsgs = mc.receivedMsgs[:0]
	}
	if mc.stateChanges != nil {
		mc.stateChanges = mc.stateChanges[:0]
	}
}

// SimulateDisconnect forces a disconnection with optional error
func (mc *MockClient[T]) SimulateDisconnect(err error) {
	if err == nil {
		err = axon.ErrConnectionClosed
	}

	mc.transitionState(axon.StateConnected, axon.StateDisconnected, err, 0)

	if mc.onDisconnect != nil {
		mc.onDisconnect(err)
	}
}

// SimulateReconnecting transitions to reconnecting state
func (mc *MockClient[T]) SimulateReconnecting(attempt int) {
	mc.transitionState(mc.State(), axon.StateReconnecting, nil, attempt)
}

// ForceState forces a specific state (use with caution in tests)
func (mc *MockClient[T]) ForceState(state axon.ConnectionState) {
	old := axon.ConnectionState(mc.state.Swap(int32(state)))

	if old != state {
		change := axon.StateChange{
			From:      old,
			To:        state,
			Time:      time.Now(),
			SessionID: mc.SessionID(),
		}

		if mc.stateChanges != nil {
			mc.mu.Lock()
			mc.stateChanges = append(mc.stateChanges, change)
			mc.mu.Unlock()
		}

		if mc.onStateChange != nil {
			mc.onStateChange(change)
		}
	}
}

// InjectConnectError injects an error for the next Connect operation
func (mc *MockClient[T]) InjectConnectError(err error) {
	if err == nil {
		mc.connectErr.Store((*errorHolder)(nil))
	} else {
		mc.connectErr.Store(&errorHolder{err: err})
	}
}

// InjectReadError injects an error for the next Read operation
func (mc *MockClient[T]) InjectReadError(err error) {
	if err == nil {
		mc.readErr.Store((*errorHolder)(nil))
	} else {
		mc.readErr.Store(&errorHolder{err: err})
	}
}

// InjectWriteError injects an error for the next Write operation
func (mc *MockClient[T]) InjectWriteError(err error) {
	if err == nil {
		mc.writeErr.Store((*errorHolder)(nil))
	} else {
		mc.writeErr.Store(&errorHolder{err: err})
	}
}

// ClearInjectedErrors clears all injected errors
func (mc *MockClient[T]) ClearInjectedErrors() {
	mc.connectErr.Store((*errorHolder)(nil))
	mc.readErr.Store((*errorHolder)(nil))
	mc.writeErr.Store((*errorHolder)(nil))
}
