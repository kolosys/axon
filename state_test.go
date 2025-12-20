package axon_test

import (
	"sync"
	"testing"
	"time"

	"github.com/kolosys/axon"
)

func TestConnectionState_String(t *testing.T) {
	tests := []struct {
		state    axon.ConnectionState
		expected string
	}{
		{axon.StateDisconnected, "disconnected"},
		{axon.StateConnecting, "connecting"},
		{axon.StateConnected, "connected"},
		{axon.StateReconnecting, "reconnecting"},
		{axon.StateClosing, "closing"},
		{axon.StateClosed, "closed"},
		{axon.ConnectionState(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConnectionState_IsActive(t *testing.T) {
	tests := []struct {
		state    axon.ConnectionState
		expected bool
	}{
		{axon.StateDisconnected, false},
		{axon.StateConnecting, true},
		{axon.StateConnected, true},
		{axon.StateReconnecting, true},
		{axon.StateClosing, false},
		{axon.StateClosed, false},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			if got := tt.state.IsActive(); got != tt.expected {
				t.Errorf("IsActive() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConnectionState_CanReconnect(t *testing.T) {
	tests := []struct {
		state    axon.ConnectionState
		expected bool
	}{
		{axon.StateDisconnected, true},
		{axon.StateConnecting, false},
		{axon.StateConnected, false},
		{axon.StateReconnecting, true},
		{axon.StateClosing, false},
		{axon.StateClosed, false},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			if got := tt.state.CanReconnect(); got != tt.expected {
				t.Errorf("CanReconnect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStateChange(t *testing.T) {
	change := axon.StateChange{
		From:      axon.StateDisconnected,
		To:        axon.StateConnecting,
		Time:      time.Now(),
		Err:       nil,
		Attempt:   0,
		SessionID: "test-session",
	}

	if change.From != axon.StateDisconnected {
		t.Errorf("From = %v, want %v", change.From, axon.StateDisconnected)
	}
	if change.To != axon.StateConnecting {
		t.Errorf("To = %v, want %v", change.To, axon.StateConnecting)
	}
	if change.SessionID != "test-session" {
		t.Errorf("SessionID = %v, want %v", change.SessionID, "test-session")
	}
}

func TestStateChange_ConcurrentAccess(t *testing.T) {
	// Test that state changes can be handled concurrently
	var wg sync.WaitGroup
	count := 100

	changes := make(chan axon.StateChange, count)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			changes <- axon.StateChange{
				From:    axon.StateDisconnected,
				To:      axon.StateConnecting,
				Time:    time.Now(),
				Attempt: i,
			}
		}(i)
	}

	wg.Wait()
	close(changes)

	received := 0
	for range changes {
		received++
	}

	if received != count {
		t.Errorf("received %d changes, want %d", received, count)
	}
}
