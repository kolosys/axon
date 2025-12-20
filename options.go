package axon

import (
	"net/http"
	"time"
)

// UpgradeOptions configures WebSocket connection upgrade options
type UpgradeOptions struct {
	// ReadBufferSize sets the size of the read buffer in bytes.
	// Default is 4096 bytes.
	ReadBufferSize int

	// WriteBufferSize sets the size of the write buffer in bytes.
	// Default is 4096 bytes.
	WriteBufferSize int

	// MaxFrameSize sets the maximum frame size in bytes.
	// Frames exceeding this size will result in ErrFrameTooLarge.
	// Default is 4096 bytes.
	MaxFrameSize int

	// MaxMessageSize sets the maximum message size in bytes.
	// Messages exceeding this size will result in ErrMessageTooLarge.
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
	// Default is 0 (disabled).
	PingInterval time.Duration

	// PongTimeout sets the timeout for waiting for a pong response.
	// If zero, pong timeout is disabled.
	// Default is 0 (disabled).
	PongTimeout time.Duration

	// CheckOrigin sets a function to validate the origin header.
	// If nil, all origins are allowed.
	// Default is nil (all origins allowed).
	CheckOrigin func(r *http.Request) bool

	// Subprotocols sets the list of supported subprotocols.
	// The client's requested subprotocol must match one of these.
	// Default is nil (no subprotocols).
	Subprotocols []string

	// Compression enables per-message compression (RFC 7692).
	// Default is false (disabled).
	Compression bool
}
