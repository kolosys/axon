package mock

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kolosys/axon"
)

// errHolder wraps an error for use with atomic.Value
type errHolder struct {
	err error
}

// MockConn provides an in-memory mock of axon.Conn[T] for testing
type MockConn[T any] struct {
	readCh  chan T
	writeCh chan T

	closed      atomic.Bool
	closeCode   atomic.Int32
	closeReason atomic.Value

	readErr  atomic.Value
	writeErr atomic.Value
	closeErr atomic.Value

	mu          sync.RWMutex
	writtenMsgs []T
	readMsgs    []T

	readDelay  time.Duration
	writeDelay time.Duration

	closeOnce sync.Once
}

// MockConnOptions configures MockConn behavior
type MockConnOptions struct {
	BufferSize     int
	ReadDelay      time.Duration
	WriteDelay     time.Duration
	RecordMessages bool
}

// NewMockConn creates a new mock connection
func NewMockConn[T any](opts *MockConnOptions) *MockConn[T] {
	if opts == nil {
		opts = &MockConnOptions{
			BufferSize:     100,
			RecordMessages: true,
		}
	}

	if opts.BufferSize <= 0 {
		opts.BufferSize = 100
	}

	mc := &MockConn[T]{
		readCh:     make(chan T, opts.BufferSize),
		writeCh:    make(chan T, opts.BufferSize),
		readDelay:  opts.ReadDelay,
		writeDelay: opts.WriteDelay,
	}

	mc.closeReason.Store("")

	if opts.RecordMessages {
		mc.writtenMsgs = make([]T, 0, 50)
		mc.readMsgs = make([]T, 0, 50)
	}

	return mc
}

// Read reads a message from the mock connection
func (mc *MockConn[T]) Read(ctx context.Context) (T, error) {
	var zero T

	if mc.closed.Load() {
		return zero, axon.ErrConnectionClosed
	}

	if eh := mc.readErr.Load(); eh != nil {
		if holder := eh.(*errHolder); holder != nil && holder.err != nil {
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

		if mc.writtenMsgs != nil {
			mc.mu.Lock()
			mc.readMsgs = append(mc.readMsgs, msg)
			mc.mu.Unlock()
		}

		return msg, nil
	}
}

// Write writes a message to the mock connection
func (mc *MockConn[T]) Write(ctx context.Context, msg T) error {
	if mc.closed.Load() {
		return axon.ErrConnectionClosed
	}

	if eh := mc.writeErr.Load(); eh != nil {
		if holder := eh.(*errHolder); holder != nil && holder.err != nil {
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

// Close closes the mock connection
func (mc *MockConn[T]) Close(code int, reason string) error {
	var closeErr error

	mc.closeOnce.Do(func() {
		mc.closed.Store(true)
		mc.closeCode.Store(int32(code))
		mc.closeReason.Store(reason)

		if eh := mc.closeErr.Load(); eh != nil {
			if holder := eh.(*errHolder); holder != nil && holder.err != nil {
				closeErr = holder.err
				return
			}
		}

		close(mc.readCh)
		close(mc.writeCh)
	})

	return closeErr
}

// IsClosed returns true if the connection is closed
func (mc *MockConn[T]) IsClosed() bool {
	return mc.closed.Load()
}

// CloseCode returns the close code
func (mc *MockConn[T]) CloseCode() int {
	return int(mc.closeCode.Load())
}

// CloseReason returns the close reason
func (mc *MockConn[T]) CloseReason() string {
	v := mc.closeReason.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

// InjectMessage simulates receiving a message (test helper)
func (mc *MockConn[T]) InjectMessage(msg T) error {
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
func (mc *MockConn[T]) DrainWritten() []T {
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
func (mc *MockConn[T]) WrittenMessages() []T {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.writtenMsgs == nil {
		return nil
	}

	msgs := make([]T, len(mc.writtenMsgs))
	copy(msgs, mc.writtenMsgs)
	return msgs
}

// ReadMessages returns recorded read messages
func (mc *MockConn[T]) ReadMessages() []T {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.readMsgs == nil {
		return nil
	}

	msgs := make([]T, len(mc.readMsgs))
	copy(msgs, mc.readMsgs)
	return msgs
}

// ClearRecorded clears all recorded messages
func (mc *MockConn[T]) ClearRecorded() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.writtenMsgs != nil {
		mc.writtenMsgs = mc.writtenMsgs[:0]
	}
	if mc.readMsgs != nil {
		mc.readMsgs = mc.readMsgs[:0]
	}
}

// InjectReadError injects an error for the next Read operation
func (mc *MockConn[T]) InjectReadError(err error) {
	if err == nil {
		mc.readErr.Store((*errHolder)(nil))
	} else {
		mc.readErr.Store(&errHolder{err: err})
	}
}

// InjectWriteError injects an error for the next Write operation
func (mc *MockConn[T]) InjectWriteError(err error) {
	if err == nil {
		mc.writeErr.Store((*errHolder)(nil))
	} else {
		mc.writeErr.Store(&errHolder{err: err})
	}
}

// InjectCloseError injects an error for the next Close operation
func (mc *MockConn[T]) InjectCloseError(err error) {
	if err == nil {
		mc.closeErr.Store((*errHolder)(nil))
	} else {
		mc.closeErr.Store(&errHolder{err: err})
	}
}

// ClearInjectedErrors clears all injected errors
func (mc *MockConn[T]) ClearInjectedErrors() {
	mc.readErr.Store((*errHolder)(nil))
	mc.writeErr.Store((*errHolder)(nil))
	mc.closeErr.Store((*errHolder)(nil))
}
