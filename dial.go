package axon

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DialOptions configures WebSocket client dial options
type DialOptions struct {
	// HandshakeTimeout is the maximum time to wait for the handshake to complete.
	// Default is 30 seconds.
	HandshakeTimeout time.Duration

	// ReadBufferSize sets the size of the read buffer in bytes.
	// Default is 4096 bytes.
	ReadBufferSize int

	// WriteBufferSize sets the size of the write buffer in bytes.
	// Default is 4096 bytes.
	WriteBufferSize int

	// MaxFrameSize sets the maximum frame size in bytes.
	// Default is 4096 bytes.
	MaxFrameSize int

	// MaxMessageSize sets the maximum message size in bytes.
	// Default is 1048576 bytes (1MB).
	MaxMessageSize int

	// ReadDeadline sets the read deadline for connections.
	// Default is no deadline.
	ReadDeadline time.Duration

	// WriteDeadline sets the write deadline for connections.
	// Default is no deadline.
	WriteDeadline time.Duration

	// PingInterval sets the interval for sending ping frames.
	// If zero, pings are disabled.
	PingInterval time.Duration

	// PongTimeout sets the timeout for waiting for a pong response.
	// If zero, pong timeout is disabled.
	PongTimeout time.Duration

	// Subprotocols sets the list of supported subprotocols.
	// Default is nil (no subprotocols).
	Subprotocols []string

	// Compression enables per-message compression (RFC 7692).
	// Default is false (disabled).
	Compression bool

	// CompressionThreshold sets the minimum message size to compress.
	// Messages smaller than this will not be compressed.
	// Default is 256 bytes.
	CompressionThreshold int

	// Headers sets additional HTTP headers for the handshake request.
	Headers http.Header

	// TLSConfig specifies the TLS configuration to use for wss:// connections.
	// If nil, the default configuration is used.
	TLSConfig *tls.Config

	// NetDialer specifies the dialer to use for creating the network connection.
	// If nil, a default dialer is used.
	NetDialer *net.Dialer
}

// Dialer is a WebSocket client dialer
type Dialer struct {
	opts *DialOptions
}

// NewDialer creates a new Dialer with the given options
func NewDialer(opts *DialOptions) *Dialer {
	if opts == nil {
		opts = &DialOptions{}
	}
	return &Dialer{opts: opts}
}

// Dial connects to a WebSocket server at the given URL.
// The URL must have a ws:// or wss:// scheme.
func Dial[T any](ctx context.Context, rawURL string, opts *DialOptions) (*Conn[T], error) {
	d := NewDialer(opts)
	return DialWithDialer[T](ctx, d, rawURL)
}

// DialWithDialer connects to a WebSocket server using the provided dialer
func DialWithDialer[T any](ctx context.Context, d *Dialer, rawURL string) (*Conn[T], error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("axon: invalid URL: %w", err)
	}

	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	default:
		return nil, fmt.Errorf("axon: unsupported scheme: %s", u.Scheme)
	}

	opts := d.opts
	if opts == nil {
		opts = &DialOptions{}
	}

	// Apply defaults
	handshakeTimeout := opts.HandshakeTimeout
	if handshakeTimeout == 0 {
		handshakeTimeout = 30 * time.Second
	}

	readBufferSize := opts.ReadBufferSize
	if readBufferSize <= 0 {
		readBufferSize = 4096
	}

	writeBufferSize := opts.WriteBufferSize
	if writeBufferSize <= 0 {
		writeBufferSize = 4096
	}

	maxFrameSize := opts.MaxFrameSize
	if maxFrameSize <= 0 {
		maxFrameSize = 4096
	}

	maxMessageSize := opts.MaxMessageSize
	if maxMessageSize <= 0 {
		maxMessageSize = 1048576
	}

	compressionThreshold := opts.CompressionThreshold
	if compressionThreshold <= 0 {
		compressionThreshold = 256
	}

	// Create context with handshake timeout
	dialCtx, dialCancel := context.WithTimeout(ctx, handshakeTimeout)
	defer dialCancel()

	// Establish TCP connection
	conn, err := d.dial(dialCtx, u)
	if err != nil {
		return nil, fmt.Errorf("axon: dial failed: %w", err)
	}

	// Generate WebSocket key
	key, err := generateWebSocketKey()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("axon: failed to generate key: %w", err)
	}

	// Build handshake request
	host := u.Host
	if u.Port() == "" {
		if u.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	path := u.Path
	if path == "" {
		path = "/"
	}
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	// Write handshake request
	var buf strings.Builder
	buf.WriteString("GET ")
	buf.WriteString(path)
	buf.WriteString(" HTTP/1.1\r\n")
	buf.WriteString("Host: ")
	buf.WriteString(host)
	buf.WriteString("\r\n")
	buf.WriteString("Upgrade: websocket\r\n")
	buf.WriteString("Connection: Upgrade\r\n")
	buf.WriteString("Sec-WebSocket-Key: ")
	buf.WriteString(key)
	buf.WriteString("\r\n")
	buf.WriteString("Sec-WebSocket-Version: 13\r\n")

	// Add subprotocols if specified
	if len(opts.Subprotocols) > 0 {
		buf.WriteString("Sec-WebSocket-Protocol: ")
		buf.WriteString(strings.Join(opts.Subprotocols, ", "))
		buf.WriteString("\r\n")
	}

	// Add compression extension if requested
	if opts.Compression {
		buf.WriteString("Sec-WebSocket-Extensions: permessage-deflate; client_max_window_bits\r\n")
	}

	// Add custom headers
	for key, values := range opts.Headers {
		for _, value := range values {
			buf.WriteString(key)
			buf.WriteString(": ")
			buf.WriteString(value)
			buf.WriteString("\r\n")
		}
	}

	buf.WriteString("\r\n")

	// Set write deadline for handshake
	if err := conn.SetWriteDeadline(time.Now().Add(handshakeTimeout)); err != nil {
		conn.Close()
		return nil, fmt.Errorf("axon: failed to set write deadline: %w", err)
	}

	if _, err := io.WriteString(conn, buf.String()); err != nil {
		conn.Close()
		return nil, fmt.Errorf("axon: failed to write handshake: %w", err)
	}

	// Set read deadline for response
	if err := conn.SetReadDeadline(time.Now().Add(handshakeTimeout)); err != nil {
		conn.Close()
		return nil, fmt.Errorf("axon: failed to set read deadline: %w", err)
	}

	// Read handshake response
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("axon: failed to read response: %w", err)
	}
	defer resp.Body.Close()

	// Validate response
	if resp.StatusCode != http.StatusSwitchingProtocols {
		conn.Close()
		return nil, fmt.Errorf("axon: unexpected status code: %d", resp.StatusCode)
	}

	if !strings.EqualFold(resp.Header.Get("Upgrade"), "websocket") {
		conn.Close()
		return nil, ErrInvalidHandshake
	}

	if !strings.Contains(strings.ToLower(resp.Header.Get("Connection")), "upgrade") {
		conn.Close()
		return nil, ErrInvalidHandshake
	}

	// Validate accept key
	expectedAccept := computeAcceptKey(key)
	if resp.Header.Get("Sec-WebSocket-Accept") != expectedAccept {
		conn.Close()
		return nil, ErrInvalidHandshake
	}

	// Check if compression was accepted
	compressionEnabled := false
	if opts.Compression {
		extensions := resp.Header.Get("Sec-WebSocket-Extensions")
		if strings.Contains(extensions, "permessage-deflate") {
			compressionEnabled = true
		}
	}

	// Clear deadlines for normal operation
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		conn.Close()
		return nil, fmt.Errorf("axon: failed to clear read deadline: %w", err)
	}
	if err := conn.SetWriteDeadline(time.Time{}); err != nil {
		conn.Close()
		return nil, fmt.Errorf("axon: failed to clear write deadline: %w", err)
	}

	// Create upgrader for connection configuration
	upgrader := &Upgrader{
		readBufferSize:    readBufferSize,
		writeBufferSize:   writeBufferSize,
		maxFrameSize:      maxFrameSize,
		maxMessageSize:    maxMessageSize,
		readDeadline:      opts.ReadDeadline,
		writeDeadline:     opts.WriteDeadline,
		pingInterval:      opts.PingInterval,
		pongTimeout:       opts.PongTimeout,
		enableCompression: compressionEnabled,
	}

	// Get pooled buffers and readers/writers
	readBuf := getBuffer()
	writeBuf := getBuffer()
	wsReader := getReader(conn)
	wsWriter := getWriter(conn)

	// Create WebSocket connection
	wsConn := &Conn[T]{
		conn:          conn,
		reader:        wsReader,
		writer:        wsWriter,
		readBuf:       readBuf,
		writeBuf:      writeBuf,
		upgrader:      upgrader,
		readDeadline:  opts.ReadDeadline,
		writeDeadline: opts.WriteDeadline,
		pingInterval:  opts.PingInterval,
		pongTimeout:   opts.PongTimeout,
		isClient:      true,
	}

	// Initialize compression if enabled
	if compressionEnabled {
		wsConn.compression = newCompressionManager(compressionThreshold)
	}

	// Start ping loop if configured
	if opts.PingInterval > 0 {
		wsConn.startPingLoop()
	}

	return wsConn, nil
}

// dial establishes a TCP connection to the server
func (d *Dialer) dial(ctx context.Context, u *url.URL) (net.Conn, error) {
	opts := d.opts
	if opts == nil {
		opts = &DialOptions{}
	}

	netDialer := opts.NetDialer
	if netDialer == nil {
		netDialer = &net.Dialer{}
	}

	host := u.Host
	if u.Port() == "" {
		if u.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	if u.Scheme == "https" {
		// TLS connection
		tlsConfig := opts.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}
		if tlsConfig.ServerName == "" {
			tlsConfig = tlsConfig.Clone()
			tlsConfig.ServerName = u.Hostname()
		}

		conn, err := tls.DialWithDialer(netDialer, "tcp", host, tlsConfig)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	// Plain TCP connection
	return netDialer.DialContext(ctx, "tcp", host)
}

// generateWebSocketKey generates a random 16-byte base64-encoded key
func generateWebSocketKey() (string, error) {
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// computeAcceptKey computes the Sec-WebSocket-Accept value for validation
func computeClientAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	h.Write([]byte(websocketGUID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
