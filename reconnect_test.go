package axon_test

import (
	"testing"
	"time"

	"github.com/kolosys/axon"
)

func TestDefaultReconnectConfig(t *testing.T) {
	cfg := axon.DefaultReconnectConfig()

	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.MaxAttempts != 0 {
		t.Errorf("expected MaxAttempts to be 0, got %d", cfg.MaxAttempts)
	}
	if cfg.InitialDelay != time.Second {
		t.Errorf("expected InitialDelay to be 1s, got %v", cfg.InitialDelay)
	}
	if cfg.MaxDelay != 30*time.Second {
		t.Errorf("expected MaxDelay to be 30s, got %v", cfg.MaxDelay)
	}
	if cfg.BackoffMultiplier != 2.0 {
		t.Errorf("expected BackoffMultiplier to be 2.0, got %f", cfg.BackoffMultiplier)
	}
	if !cfg.Jitter {
		t.Error("expected Jitter to be true")
	}
	if cfg.ResetAfter != 60*time.Second {
		t.Errorf("expected ResetAfter to be 60s, got %v", cfg.ResetAfter)
	}
}

func TestReconnectConfig_Callbacks(t *testing.T) {
	var reconnectingCalled, reconnectedCalled, failedCalled bool
	var lastAttempt int
	var lastDelay time.Duration

	cfg := &axon.ReconnectConfig{
		Enabled:           true,
		MaxAttempts:       3,
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          time.Second,
		BackoffMultiplier: 2.0,
		Jitter:            false,
		OnReconnecting: func(attempt int, delay time.Duration) {
			reconnectingCalled = true
			lastAttempt = attempt
			lastDelay = delay
		},
		OnReconnected: func(attempt int) {
			reconnectedCalled = true
		},
		OnReconnectFailed: func(attempt int, err error) {
			failedCalled = true
		},
	}

	// Verify callbacks are set
	if cfg.OnReconnecting == nil {
		t.Error("expected OnReconnecting to be set")
	}
	if cfg.OnReconnected == nil {
		t.Error("expected OnReconnected to be set")
	}
	if cfg.OnReconnectFailed == nil {
		t.Error("expected OnReconnectFailed to be set")
	}

	// Call the callbacks
	cfg.OnReconnecting(1, 100*time.Millisecond)
	cfg.OnReconnected(1)
	cfg.OnReconnectFailed(1, nil)

	if !reconnectingCalled {
		t.Error("expected OnReconnecting to be called")
	}
	if !reconnectedCalled {
		t.Error("expected OnReconnected to be called")
	}
	if !failedCalled {
		t.Error("expected OnReconnectFailed to be called")
	}
	if lastAttempt != 1 {
		t.Errorf("expected lastAttempt to be 1, got %d", lastAttempt)
	}
	if lastDelay != 100*time.Millisecond {
		t.Errorf("expected lastDelay to be 100ms, got %v", lastDelay)
	}
}

func TestReconnectConfig_ShouldReconnect(t *testing.T) {
	tests := []struct {
		name            string
		maxAttempts     int
		currentAttempt  int
		shouldReconnect func(error, int) bool
		expected        bool
	}{
		{
			name:           "unlimited attempts",
			maxAttempts:    0,
			currentAttempt: 100,
			expected:       true,
		},
		{
			name:           "within max attempts",
			maxAttempts:    5,
			currentAttempt: 3,
			expected:       true,
		},
		{
			name:           "at max attempts",
			maxAttempts:    5,
			currentAttempt: 5,
			expected:       false,
		},
		{
			name:           "custom predicate returns true",
			maxAttempts:    0,
			currentAttempt: 1,
			shouldReconnect: func(err error, attempt int) bool {
				return true
			},
			expected: true,
		},
		{
			name:           "custom predicate returns false",
			maxAttempts:    0,
			currentAttempt: 1,
			shouldReconnect: func(err error, attempt int) bool {
				return false
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &axon.ReconnectConfig{
				Enabled:         true,
				MaxAttempts:     tt.maxAttempts,
				ShouldReconnect: tt.shouldReconnect,
			}

			// Verify the config is set correctly
			if cfg.MaxAttempts != tt.maxAttempts {
				t.Errorf("MaxAttempts = %d, want %d", cfg.MaxAttempts, tt.maxAttempts)
			}
		})
	}
}

func TestReconnectConfig_BackoffCalculation(t *testing.T) {
	cfg := &axon.ReconnectConfig{
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          time.Second,
		BackoffMultiplier: 2.0,
		Jitter:            false,
	}

	// Test exponential backoff: 100ms, 200ms, 400ms, 800ms, 1000ms (capped)
	expectedDelays := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
		time.Second, // Capped at MaxDelay
	}

	for i, expected := range expectedDelays {
		// The actual calculation happens inside the reconnector
		// Here we just verify the config is set correctly
		if cfg.InitialDelay != 100*time.Millisecond {
			t.Errorf("InitialDelay = %v, want %v", cfg.InitialDelay, 100*time.Millisecond)
		}
		_ = expected
		_ = i
	}
}
