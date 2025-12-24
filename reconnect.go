package axon

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// ReconnectConfig configures automatic reconnection behavior
type ReconnectConfig struct {
	// Enabled determines if automatic reconnection is active
	Enabled bool

	// MaxAttempts is the maximum number of reconnection attempts (0 = unlimited)
	MaxAttempts int

	// InitialDelay is the initial delay before the first reconnection attempt
	// Default: 1 second
	InitialDelay time.Duration

	// MaxDelay is the maximum delay between reconnection attempts
	// Default: 30 seconds
	MaxDelay time.Duration

	// BackoffMultiplier is the multiplier for exponential backoff
	// Default: 2.0
	BackoffMultiplier float64

	// Jitter adds randomness to delay to prevent thundering herd
	// Default: true
	Jitter bool

	// ResetAfter resets the attempt counter after being connected for this duration
	// Default: 60 seconds
	ResetAfter time.Duration

	// OnReconnecting is called when a reconnection attempt starts
	OnReconnecting func(attempt int, delay time.Duration)

	// OnReconnected is called when reconnection succeeds
	OnReconnected func(attempt int)

	// OnReconnectFailed is called when reconnection fails
	OnReconnectFailed func(attempt int, err error)

	// ShouldReconnect is called to determine if reconnection should be attempted
	// If nil, always attempts to reconnect (within MaxAttempts)
	ShouldReconnect func(err error, attempt int) bool
}

// DefaultReconnectConfig returns a default reconnection configuration
func DefaultReconnectConfig() *ReconnectConfig {
	return &ReconnectConfig{
		Enabled:           true,
		MaxAttempts:       0, // Unlimited
		InitialDelay:      time.Second,
		MaxDelay:          30 * time.Second,
		BackoffMultiplier: 2.0,
		Jitter:            true,
		ResetAfter:        60 * time.Second,
	}
}

// reconnector manages automatic reconnection
type reconnector struct {
	config      *ReconnectConfig
	attempts    int
	lastConnect time.Time
	rand        *rand.Rand
}

// newReconnector creates a new reconnector with the given configuration
func newReconnector(config *ReconnectConfig) *reconnector {
	if config == nil {
		config = DefaultReconnectConfig()
	}

	if config.InitialDelay <= 0 {
		config.InitialDelay = time.Second
	}
	if config.MaxDelay <= 0 {
		config.MaxDelay = 30 * time.Second
	}
	if config.BackoffMultiplier <= 0 {
		config.BackoffMultiplier = 2.0
	}
	if config.ResetAfter <= 0 {
		config.ResetAfter = 60 * time.Second
	}

	return &reconnector{
		config: config,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// shouldReconnect determines if a reconnection should be attempted
func (r *reconnector) shouldReconnect(err error) bool {
	if !r.config.Enabled {
		return false
	}

	if r.config.MaxAttempts > 0 && r.attempts >= r.config.MaxAttempts {
		return false
	}

	if r.config.ShouldReconnect != nil {
		return r.config.ShouldReconnect(err, r.attempts)
	}

	if closeErr := AsCloseError(err); closeErr != nil {
		return closeErr.IsRecoverable()
	}

	return true
}

// nextDelay calculates the next reconnection delay using exponential backoff
func (r *reconnector) nextDelay() time.Duration {
	delay := float64(r.config.InitialDelay) * math.Pow(r.config.BackoffMultiplier, float64(r.attempts))

	// Cap at max delay
	if delay > float64(r.config.MaxDelay) {
		delay = float64(r.config.MaxDelay)
	}

	// Add jitter if enabled (±25%)
	if r.config.Jitter {
		jitterRange := delay * 0.25
		jitter := (r.rand.Float64() * 2 * jitterRange) - jitterRange
		delay += jitter
	}

	return time.Duration(delay)
}

// attempt performs a single reconnection attempt
func (r *reconnector) attempt(ctx context.Context, dialFn func(context.Context) error) error {
	r.attempts++
	delay := r.nextDelay()

	if r.config.OnReconnecting != nil {
		r.config.OnReconnecting(r.attempts, delay)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
	}

	err := dialFn(ctx)
	if err != nil {
		if r.config.OnReconnectFailed != nil {
			r.config.OnReconnectFailed(r.attempts, err)
		}
		return err
	}

	if r.config.OnReconnected != nil {
		r.config.OnReconnected(r.attempts) // Success
	}

	r.lastConnect = time.Now()
	return nil
}

// reset resets the attempt counter
func (r *reconnector) reset() {
	r.attempts = 0
}

// maybeReset resets the attempt counter if connected long enough
func (r *reconnector) maybeReset() {
	if !r.lastConnect.IsZero() && time.Since(r.lastConnect) >= r.config.ResetAfter {
		r.reset()
	}
}

// reconnectLoop continuously attempts to reconnect until successful or cancelled
func (r *reconnector) reconnectLoop(ctx context.Context, dialFn func(context.Context) error) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !r.shouldReconnect(nil) {
			return ErrReconnectFailed
		}

		err := r.attempt(ctx, dialFn)
		if err == nil {
			return nil // Success
		}

		if !r.shouldReconnect(err) {
			return ErrReconnectFailed
		}
	}
}
