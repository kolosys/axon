package axon

import "errors"

// Sentinel errors for WebSocket protocol violations and state errors
var (
	// ErrInvalidHandshake indicates the WebSocket handshake failed
	ErrInvalidHandshake = errors.New("axon: invalid websocket handshake")

	// ErrUpgradeRequired indicates the request is not a valid WebSocket upgrade request
	ErrUpgradeRequired = errors.New("axon: upgrade required")

	// ErrInvalidOrigin indicates the origin header is not allowed
	ErrInvalidOrigin = errors.New("axon: invalid origin")

	// ErrInvalidSubprotocol indicates the requested subprotocol is not supported
	ErrInvalidSubprotocol = errors.New("axon: invalid subprotocol")

	// ErrConnectionClosed indicates the connection has been closed
	ErrConnectionClosed = errors.New("axon: connection closed")

	// ErrFrameTooLarge indicates a frame exceeds the maximum allowed size
	ErrFrameTooLarge = errors.New("axon: frame too large")

	// ErrMessageTooLarge indicates a message exceeds the maximum allowed size
	ErrMessageTooLarge = errors.New("axon: message too large")

	// ErrInvalidFrame indicates a frame violates the WebSocket protocol
	ErrInvalidFrame = errors.New("axon: invalid frame")

	// ErrInvalidMask indicates frame masking is invalid
	ErrInvalidMask = errors.New("axon: invalid mask")

	// ErrUnsupportedFrameType indicates the frame type is not supported
	ErrUnsupportedFrameType = errors.New("axon: unsupported frame type")

	// ErrFragmentedControlFrame indicates a control frame is fragmented (not allowed)
	ErrFragmentedControlFrame = errors.New("axon: fragmented control frame")

	// ErrInvalidCloseCode indicates an invalid close code was used
	ErrInvalidCloseCode = errors.New("axon: invalid close code")

	// ErrReadDeadlineExceeded indicates a read operation exceeded its deadline
	ErrReadDeadlineExceeded = errors.New("axon: read deadline exceeded")

	// ErrWriteDeadlineExceeded indicates a write operation exceeded its deadline
	ErrWriteDeadlineExceeded = errors.New("axon: write deadline exceeded")

	// ErrContextCanceled indicates the context was canceled
	ErrContextCanceled = errors.New("axon: context canceled")

	// ErrSerializationFailed indicates message serialization failed
	ErrSerializationFailed = errors.New("axon: serialization failed")

	// ErrDeserializationFailed indicates message deserialization failed
	ErrDeserializationFailed = errors.New("axon: deserialization failed")

	// ErrReconnectFailed indicates reconnection attempts have been exhausted
	ErrReconnectFailed = errors.New("axon: reconnection failed")

	// ErrQueueFull indicates the message queue is full
	ErrQueueFull = errors.New("axon: message queue full")

	// ErrQueueClosed indicates the message queue has been closed
	ErrQueueClosed = errors.New("axon: message queue closed")

	// ErrQueueTimeout indicates a queued message has expired
	ErrQueueTimeout = errors.New("axon: queued message timeout")

	// ErrQueueCleared indicates the queue was cleared before message was sent
	ErrQueueCleared = errors.New("axon: message queue cleared")

	// ErrCompressionFailed indicates compression or decompression failed
	ErrCompressionFailed = errors.New("axon: compression failed")

	// ErrInvalidState indicates an invalid state transition was attempted
	ErrInvalidState = errors.New("axon: invalid state transition")

	// ErrClientClosed indicates the client has been closed
	ErrClientClosed = errors.New("axon: client closed")
)
