package axon_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kolosys/axon"
)

func TestDial_InvalidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"empty URL", ""},
		{"invalid scheme", "http://localhost:8080"},
		{"invalid scheme ftp", "ftp://localhost:8080"},
		{"missing scheme", "localhost:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			_, err := axon.Dial[string](ctx, tt.url, nil)
			if err == nil {
				t.Error("expected error for invalid URL")
			}
		})
	}
}

func TestDial_WithOptions(t *testing.T) {
	// Create a test WebSocket server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := axon.Upgrade[string](w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close(1000, "done")

		// Just keep the connection open briefly
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := &axon.DialOptions{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   8192,
		WriteBufferSize:  8192,
		MaxFrameSize:     16384,
		MaxMessageSize:   2097152,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := axon.Dial[string](ctx, wsURL, opts)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	defer conn.Close(1000, "done")

	// Verify connection is not nil
	if conn == nil {
		t.Error("expected connection to be non-nil")
	}
}

func TestDial_WithSubprotocols(t *testing.T) {
	// Create a test WebSocket server that accepts a subprotocol
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := &axon.UpgradeOptions{
			Subprotocols: []string{"chat", "json"},
		}
		conn, err := axon.Upgrade[string](w, r, opts)
		if err != nil {
			return
		}
		defer conn.Close(1000, "done")
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := &axon.DialOptions{
		HandshakeTimeout: 5 * time.Second,
		Subprotocols:     []string{"chat"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := axon.Dial[string](ctx, wsURL, opts)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	defer conn.Close(1000, "done")
}

func TestDial_WithCompression(t *testing.T) {
	// Create a test WebSocket server with compression
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := &axon.UpgradeOptions{
			Compression: true,
		}
		conn, err := axon.Upgrade[string](w, r, opts)
		if err != nil {
			return
		}
		defer conn.Close(1000, "done")
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := &axon.DialOptions{
		HandshakeTimeout:     5 * time.Second,
		Compression:          true,
		CompressionThreshold: 256,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := axon.Dial[string](ctx, wsURL, opts)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	defer conn.Close(1000, "done")
}

func TestDial_WithCustomHeaders(t *testing.T) {
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		conn, err := axon.Upgrade[string](w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close(1000, "done")
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-token")

	opts := &axon.DialOptions{
		HandshakeTimeout: 5 * time.Second,
		Headers:          headers,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := axon.Dial[string](ctx, wsURL, opts)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	defer conn.Close(1000, "done")

	// Give the server time to process
	time.Sleep(50 * time.Millisecond)

	if receivedAuth != "Bearer test-token" {
		t.Errorf("Authorization header = %q, want %q", receivedAuth, "Bearer test-token")
	}
}

func TestDial_ConnectionRefused(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Try to connect to a port that's not listening
	_, err := axon.Dial[string](ctx, "ws://localhost:59999", nil)
	if err == nil {
		t.Error("expected error for connection refused")
	}
}

func TestDial_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := axon.Dial[string](ctx, "ws://localhost:8080", nil)
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestDialOptions_Defaults(t *testing.T) {
	opts := &axon.DialOptions{}

	// Verify zero values
	if opts.HandshakeTimeout != 0 {
		t.Errorf("HandshakeTimeout = %v, want 0", opts.HandshakeTimeout)
	}
	if opts.ReadBufferSize != 0 {
		t.Errorf("ReadBufferSize = %d, want 0", opts.ReadBufferSize)
	}
	if opts.WriteBufferSize != 0 {
		t.Errorf("WriteBufferSize = %d, want 0", opts.WriteBufferSize)
	}
}

func TestNewDialer(t *testing.T) {
	dialer := axon.NewDialer(nil)
	if dialer == nil {
		t.Error("expected dialer to be non-nil")
	}

	opts := &axon.DialOptions{
		HandshakeTimeout: 10 * time.Second,
	}
	dialer = axon.NewDialer(opts)
	if dialer == nil {
		t.Error("expected dialer to be non-nil")
	}
}

func TestDialWithDialer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := axon.Upgrade[string](w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close(1000, "done")
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	dialer := axon.NewDialer(&axon.DialOptions{
		HandshakeTimeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := axon.DialWithDialer[string](ctx, dialer, wsURL)
	if err != nil {
		t.Fatalf("DialWithDialer() error = %v", err)
	}
	defer conn.Close(1000, "done")
}
