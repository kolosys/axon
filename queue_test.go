package axon_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kolosys/axon"
)

func TestMessageQueueStats(t *testing.T) {
	stats := axon.MessageQueueStats{
		CurrentSize: 10,
		MaxSize:     100,
		Dropped:     5,
		Enqueued:    50,
		Sent:        35,
	}

	if stats.CurrentSize != 10 {
		t.Errorf("CurrentSize = %d, want 10", stats.CurrentSize)
	}
	if stats.MaxSize != 100 {
		t.Errorf("MaxSize = %d, want 100", stats.MaxSize)
	}
	if stats.Dropped != 5 {
		t.Errorf("Dropped = %d, want 5", stats.Dropped)
	}
	if stats.Enqueued != 50 {
		t.Errorf("Enqueued = %d, want 50", stats.Enqueued)
	}
	if stats.Sent != 35 {
		t.Errorf("Sent = %d, want 35", stats.Sent)
	}
}

func TestClientOptions_QueueConfig(t *testing.T) {
	opts := &axon.ClientOptions{
		QueueSize:    100,
		QueueTimeout: 30 * time.Second,
	}

	if opts.QueueSize != 100 {
		t.Errorf("QueueSize = %d, want 100", opts.QueueSize)
	}
	if opts.QueueTimeout != 30*time.Second {
		t.Errorf("QueueTimeout = %v, want 30s", opts.QueueTimeout)
	}
}

func TestDefaultClientOptions_Queue(t *testing.T) {
	opts := axon.DefaultClientOptions()

	if opts.QueueSize != 100 {
		t.Errorf("QueueSize = %d, want 100", opts.QueueSize)
	}
	if opts.QueueTimeout != 30*time.Second {
		t.Errorf("QueueTimeout = %v, want 30s", opts.QueueTimeout)
	}
}

func TestQueueErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"QueueFull", axon.ErrQueueFull},
		{"QueueClosed", axon.ErrQueueClosed},
		{"QueueTimeout", axon.ErrQueueTimeout},
		{"QueueCleared", axon.ErrQueueCleared},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("expected error to be non-nil")
			}
			if tt.err.Error() == "" {
				t.Error("expected error message to be non-empty")
			}
		})
	}
}

func TestQueueStats_ConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	count := 100

	stats := make([]axon.MessageQueueStats, count)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			stats[i] = axon.MessageQueueStats{
				CurrentSize: i,
				MaxSize:     100,
				Dropped:     int64(i),
				Enqueued:    int64(i * 2),
				Sent:        int64(i),
			}
		}(i)
	}

	wg.Wait()

	// Verify all stats were created
	for i := 0; i < count; i++ {
		if stats[i].CurrentSize != i {
			t.Errorf("stats[%d].CurrentSize = %d, want %d", i, stats[i].CurrentSize, i)
		}
	}
}

func TestClient_QueueStats(t *testing.T) {
	opts := &axon.ClientOptions{
		QueueSize:    10,
		QueueTimeout: time.Second,
	}

	client := axon.NewClient[string]("ws://localhost:8080", opts)
	defer client.Close()

	stats := client.QueueStats()

	// Initially, queue should be empty
	if stats.CurrentSize != 0 {
		t.Errorf("CurrentSize = %d, want 0", stats.CurrentSize)
	}
	if stats.MaxSize != 10 {
		t.Errorf("MaxSize = %d, want 10", stats.MaxSize)
	}
}

func TestClient_QueueDisabled(t *testing.T) {
	opts := &axon.ClientOptions{
		QueueSize: 0, // Disable queue
	}

	client := axon.NewClient[string]("ws://localhost:8080", opts)
	defer client.Close()

	stats := client.QueueStats()

	// When queue is disabled, stats should be zero
	if stats.MaxSize != 0 {
		t.Errorf("MaxSize = %d, want 0 (queue disabled)", stats.MaxSize)
	}
}

func TestClient_WriteWhenDisconnected(t *testing.T) {
	opts := &axon.ClientOptions{
		QueueSize:    10,
		QueueTimeout: 100 * time.Millisecond,
	}

	client := axon.NewClient[string]("ws://localhost:8080", opts)
	defer client.Close()

	// Try to write when disconnected - should get an error since not in reconnecting state
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := client.Write(ctx, "test message")
	if err == nil {
		t.Error("expected error when writing while disconnected")
	}
}
