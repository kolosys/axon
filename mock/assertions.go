package mock

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/kolosys/axon"
)

// AssertWritten asserts that a message was written
func AssertWritten[T any](t *testing.T, mc *MockConn[T], expected T) {
	t.Helper()

	msgs := mc.WrittenMessages()
	for _, msg := range msgs {
		if reflect.DeepEqual(msg, expected) {
			return
		}
	}

	t.Errorf("expected message %v to be written, but it was not found in %v", expected, msgs)
}

// AssertWrittenCount asserts the number of written messages
func AssertWrittenCount[T any](t *testing.T, mc *MockConn[T], expected int) {
	t.Helper()

	msgs := mc.WrittenMessages()
	if len(msgs) != expected {
		t.Errorf("expected %d written messages, got %d", expected, len(msgs))
	}
}

// AssertClientWritten asserts that a message was written by the client
func AssertClientWritten[T any](t *testing.T, mc *MockClient[T], expected T) {
	t.Helper()

	msgs := mc.WrittenMessages()
	for _, msg := range msgs {
		if reflect.DeepEqual(msg, expected) {
			return
		}
	}

	t.Errorf("expected message %v to be written, but it was not found in %v", expected, msgs)
}

// AssertClientWrittenCount asserts the number of written messages by the client
func AssertClientWrittenCount[T any](t *testing.T, mc *MockClient[T], expected int) {
	t.Helper()

	msgs := mc.WrittenMessages()
	if len(msgs) != expected {
		t.Errorf("expected %d written messages, got %d", expected, len(msgs))
	}
}

// AssertState asserts the client is in the expected state
func AssertState[T any](t *testing.T, mc *MockClient[T], expected axon.ConnectionState) {
	t.Helper()

	actual := mc.State()
	if actual != expected {
		t.Errorf("expected state %v, got %v", expected, actual)
	}
}

// AssertStateTransition asserts a state transition occurred
func AssertStateTransition[T any](t *testing.T, mc *MockClient[T], from, to axon.ConnectionState) {
	t.Helper()

	changes := mc.StateChanges()
	for _, change := range changes {
		if change.From == from && change.To == to {
			return
		}
	}

	t.Errorf("expected state transition %v -> %v, but it did not occur", from, to)
}

// AssertConnected asserts the client is connected
func AssertConnected[T any](t *testing.T, mc *MockClient[T]) {
	t.Helper()
	AssertState(t, mc, axon.StateConnected)
}

// AssertDisconnected asserts the client is disconnected
func AssertDisconnected[T any](t *testing.T, mc *MockClient[T]) {
	t.Helper()
	AssertState(t, mc, axon.StateDisconnected)
}

// AssertClosed asserts the client is closed
func AssertClosed[T any](t *testing.T, mc *MockClient[T]) {
	t.Helper()
	AssertState(t, mc, axon.StateClosed)
}

// WaitForState waits for the client to reach the expected state
func WaitForState[T any](t *testing.T, mc *MockClient[T], expected axon.ConnectionState, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if mc.State() == expected {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("timeout waiting for state %v, current state: %v", expected, mc.State())
}

// WaitForMessage waits for a message to be written
func WaitForMessage[T any](t *testing.T, mc *MockConn[T], expected T, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		msgs := mc.WrittenMessages()
		for _, msg := range msgs {
			if reflect.DeepEqual(msg, expected) {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("timeout waiting for message %v", expected)
}

// WaitForClientMessage waits for a message to be written by the client
func WaitForClientMessage[T any](t *testing.T, mc *MockClient[T], expected T, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		msgs := mc.WrittenMessages()
		for _, msg := range msgs {
			if reflect.DeepEqual(msg, expected) {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("timeout waiting for message %v", expected)
}

// CallbackRecorder records callback invocations for testing
type CallbackRecorder struct {
	mu     sync.Mutex
	called int
}

// NewCallbackRecorder creates a new callback recorder
func NewCallbackRecorder() *CallbackRecorder {
	return &CallbackRecorder{}
}

// Record records a callback invocation
func (cr *CallbackRecorder) Record() {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.called++
}

// Called returns the number of times the callback was called
func (cr *CallbackRecorder) Called() int {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	return cr.called
}

// AssertCalled asserts a callback was invoked
func (cr *CallbackRecorder) AssertCalled(t *testing.T) {
	t.Helper()
	if cr.called == 0 {
		t.Error("expected callback to be called, but it was not")
	}
}

// AssertNotCalled asserts a callback was not invoked
func (cr *CallbackRecorder) AssertNotCalled(t *testing.T) {
	t.Helper()
	if cr.called > 0 {
		t.Errorf("expected callback not to be called, but it was called %d times", cr.called)
	}
}

// AssertCallCount asserts the exact number of callback invocations
func (cr *CallbackRecorder) AssertCallCount(t *testing.T, expected int) {
	t.Helper()
	if cr.called != expected {
		t.Errorf("expected callback to be called %d times, got %d", expected, cr.called)
	}
}

// AssertStateChangeCallback asserts state changes via callback
func AssertStateChangeCallback[T any](t *testing.T, mc *MockClient[T], from, to axon.ConnectionState) {
	t.Helper()

	recorder := NewCallbackRecorder()
	mc.OnStateChange(func(change axon.StateChange) {
		if change.From == from && change.To == to {
			recorder.Record()
		}
	})

	if recorder.Called() == 0 {
		t.Errorf("expected state change callback for %v -> %v", from, to)
	}
}

// MessageMatcher is a function that matches messages for assertions
type MessageMatcher[T any] func(T) bool

// AssertWrittenMatches asserts a message matching the predicate was written
func AssertWrittenMatches[T any](t *testing.T, mc *MockConn[T], matcher MessageMatcher[T]) {
	t.Helper()

	msgs := mc.WrittenMessages()
	for _, msg := range msgs {
		if matcher(msg) {
			return
		}
	}

	t.Errorf("expected to find a message matching the predicate in written messages")
}

// AssertClientWrittenMatches asserts a message matching the predicate was written by the client
func AssertClientWrittenMatches[T any](t *testing.T, mc *MockClient[T], matcher MessageMatcher[T]) {
	t.Helper()

	msgs := mc.WrittenMessages()
	for _, msg := range msgs {
		if matcher(msg) {
			return
		}
	}

	t.Errorf("expected to find a message matching the predicate in written messages")
}

// AssertReceivedCount asserts the number of received messages
func AssertReceivedCount[T any](t *testing.T, mc *MockClient[T], expected int) {
	t.Helper()

	msgs := mc.ReceivedMessages()
	if len(msgs) != expected {
		t.Errorf("expected %d received messages, got %d", expected, len(msgs))
	}
}

// AssertStateChangeCount asserts the number of state changes
func AssertStateChangeCount[T any](t *testing.T, mc *MockClient[T], expected int) {
	t.Helper()

	changes := mc.StateChanges()
	if len(changes) != expected {
		t.Errorf("expected %d state changes, got %d", expected, len(changes))
	}
}

// AssertLastStateChange asserts the last state change matches expectations
func AssertLastStateChange[T any](t *testing.T, mc *MockClient[T], from, to axon.ConnectionState) {
	t.Helper()

	changes := mc.StateChanges()
	if len(changes) == 0 {
		t.Fatal("expected at least one state change")
	}

	last := changes[len(changes)-1]
	if last.From != from || last.To != to {
		t.Errorf("last state change = %v -> %v, want %v -> %v", last.From, last.To, from, to)
	}
}

// PanicOnError returns a function that panics on error (for use in callbacks)
func PanicOnError() func(error) {
	return func(err error) {
		if err != nil {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}
	}
}

// NoOpCallback returns a no-op callback for testing without side effects
func NoOpCallback[T any]() func(T) {
	return func(T) {}
}
