package axon_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kolosys/axon"
)

func TestNewClient(t *testing.T) {
	client := axon.NewClient[string]("ws://localhost:8080", nil)
	defer client.Close()

	if client == nil {
		t.Fatal("expected client to be non-nil")
	}

	// Check initial state
	if client.State() != axon.StateDisconnected {
		t.Errorf("State() = %v, want %v", client.State(), axon.StateDisconnected)
	}

	if client.IsConnected() {
		t.Error("expected IsConnected() to be false")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	opts := &axon.ClientOptions{
		DialOptions: axon.DialOptions{
			HandshakeTimeout: 5 * time.Second,
			ReadBufferSize:   8192,
			WriteBufferSize:  8192,
			PingInterval:     10 * time.Second,
		},
		Reconnect: &axon.ReconnectConfig{
			Enabled:      true,
			MaxAttempts:  5,
			InitialDelay: 500 * time.Millisecond,
		},
		QueueSize:    50,
		QueueTimeout: 10 * time.Second,
	}

	client := axon.NewClient[string]("ws://localhost:8080", opts)
	defer client.Close()

	if client == nil {
		t.Fatal("expected client to be non-nil")
	}

	// Check queue stats reflect the options
	stats := client.QueueStats()
	if stats.MaxSize != 50 {
		t.Errorf("QueueStats().MaxSize = %d, want 50", stats.MaxSize)
	}
}

func TestClient_State(t *testing.T) {
	client := axon.NewClient[string]("ws://localhost:8080", nil)
	defer client.Close()

	// Initial state should be disconnected
	if client.State() != axon.StateDisconnected {
		t.Errorf("State() = %v, want %v", client.State(), axon.StateDisconnected)
	}

	// Close the client
	client.Close()

	// State should be closed after Close()
	if client.State() != axon.StateClosed {
		t.Errorf("State() = %v, want %v", client.State(), axon.StateClosed)
	}
}

func TestClient_OnStateChange(t *testing.T) {
	client := axon.NewClient[string]("ws://localhost:8080", nil)
	defer client.Close()

	var changes []axon.StateChange
	var mu sync.Mutex

	client.OnStateChange(func(change axon.StateChange) {
		mu.Lock()
		changes = append(changes, change)
		mu.Unlock()
	})

	// Close should trigger state changes
	client.Close()

	// Give some time for state changes to be processed
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	// Should have at least one state change (to Closing and/or Closed)
	if len(changes) == 0 {
		t.Error("expected at least one state change")
	}
}

func TestClient_SessionID(t *testing.T) {
	client := axon.NewClient[string]("ws://localhost:8080", nil)
	defer client.Close()

	// Initial session ID should be empty
	if client.SessionID() != "" {
		t.Errorf("SessionID() = %q, want empty", client.SessionID())
	}

	// Set session ID
	client.SetSessionID("test-session-123")

	if client.SessionID() != "test-session-123" {
		t.Errorf("SessionID() = %q, want %q", client.SessionID(), "test-session-123")
	}
}

func TestClient_Callbacks(t *testing.T) {
	client := axon.NewClient[string]("ws://localhost:8080", nil)
	defer client.Close()

	var connectCalled, disconnectCalled, messageCalled bool

	client.OnConnect(func(c *axon.Client[string]) {
		connectCalled = true
	})

	client.OnDisconnect(func(c *axon.Client[string], err error) {
		disconnectCalled = true
	})

	client.OnMessage(func(msg string) {
		messageCalled = true
	})

	// Verify callbacks are registered (we can't easily test them being called
	// without a real WebSocket server)
	_ = connectCalled
	_ = disconnectCalled
	_ = messageCalled
}

func TestClient_Close(t *testing.T) {
	client := axon.NewClient[string]("ws://localhost:8080", nil)

	// First close should succeed
	err := client.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}

	// Second close should be idempotent
	err = client.Close()
	if err != nil {
		t.Errorf("second Close() error = %v, want nil", err)
	}

	// State should be closed
	if client.State() != axon.StateClosed {
		t.Errorf("State() = %v, want %v", client.State(), axon.StateClosed)
	}
}

func TestClient_Conn(t *testing.T) {
	client := axon.NewClient[string]("ws://localhost:8080", nil)
	defer client.Close()

	// Before connecting, Conn() should return nil
	if client.Conn() != nil {
		t.Error("expected Conn() to be nil before connecting")
	}
}

func TestClient_ConnectTimeout(t *testing.T) {
	opts := &axon.ClientOptions{
		DialOptions: axon.DialOptions{
			HandshakeTimeout: 100 * time.Millisecond,
		},
	}

	client := axon.NewClient[string]("ws://localhost:59999", opts) // Non-existent port
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := client.Connect(ctx)
	if err == nil {
		t.Error("expected error when connecting to non-existent server")
	}

	// State should be back to disconnected after failed connect
	if client.State() != axon.StateDisconnected {
		t.Errorf("State() = %v, want %v", client.State(), axon.StateDisconnected)
	}
}

func TestClient_ConnectWithReadLoop_InvalidURL(t *testing.T) {
	client := axon.NewClient[string]("invalid://url", nil)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := client.ConnectWithReadLoop(ctx)
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestDefaultClientOptions(t *testing.T) {
	opts := axon.DefaultClientOptions()

	if opts.HandshakeTimeout != 30*time.Second {
		t.Errorf("HandshakeTimeout = %v, want 30s", opts.HandshakeTimeout)
	}
	if opts.ReadBufferSize != 4096 {
		t.Errorf("ReadBufferSize = %d, want 4096", opts.ReadBufferSize)
	}
	if opts.WriteBufferSize != 4096 {
		t.Errorf("WriteBufferSize = %d, want 4096", opts.WriteBufferSize)
	}
	if opts.MaxFrameSize != 4096 {
		t.Errorf("MaxFrameSize = %d, want 4096", opts.MaxFrameSize)
	}
	if opts.MaxMessageSize != 1048576 {
		t.Errorf("MaxMessageSize = %d, want 1048576", opts.MaxMessageSize)
	}
	if opts.PingInterval != 30*time.Second {
		t.Errorf("PingInterval = %v, want 30s", opts.PingInterval)
	}
	if opts.PongTimeout != 10*time.Second {
		t.Errorf("PongTimeout = %v, want 10s", opts.PongTimeout)
	}
	if opts.Reconnect == nil {
		t.Error("expected Reconnect to be non-nil")
	}
	if opts.QueueSize != 100 {
		t.Errorf("QueueSize = %d, want 100", opts.QueueSize)
	}
	if opts.QueueTimeout != 30*time.Second {
		t.Errorf("QueueTimeout = %v, want 30s", opts.QueueTimeout)
	}
}

func TestClient_WithRealServer(t *testing.T) {
	// Create a test WebSocket server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := axon.Upgrade[string](w, r, nil)
		if err != nil {
			t.Logf("upgrade error: %v", err)
			return
		}
		defer conn.Close(1000, "done")

		// Echo messages back
		for {
			msg, err := conn.Read(r.Context())
			if err != nil {
				return
			}
			if err := conn.Write(r.Context(), msg); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := &axon.ClientOptions{
		DialOptions: axon.DialOptions{
			HandshakeTimeout: 5 * time.Second,
		},
		Reconnect: &axon.ReconnectConfig{
			Enabled: false, // Disable reconnection for this test
		},
		QueueSize: 10,
	}

	client := axon.NewClient[string](wsURL, opts)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to the server
	err := client.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	// Verify connected state
	if client.State() != axon.StateConnected {
		t.Errorf("State() = %v, want %v", client.State(), axon.StateConnected)
	}

	if !client.IsConnected() {
		t.Error("expected IsConnected() to be true")
	}

	// Write a message
	err = client.Write(ctx, "Hello, World!")
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Read the echo
	msg, err := client.Read(ctx)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if msg != "Hello, World!" {
		t.Errorf("Read() = %q, want %q", msg, "Hello, World!")
	}

	// Close the client
	err = client.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if client.State() != axon.StateClosed {
		t.Errorf("State() = %v, want %v", client.State(), axon.StateClosed)
	}
}
