package axon

import (
	"context"
	"sync"
	"time"
)

// Client is a WebSocket client with automatic reconnection and message queuing
type Client[T any] struct {
	// Connection management
	url    string
	opts   *ClientOptions
	conn   *Conn[T]
	connMu sync.RWMutex
	dialer *Dialer

	// State management
	state *stateManager

	// Reconnection
	reconnector *reconnector

	// Message queue
	queue *MessageQueue[T]

	// Callbacks
	onConnect    func(*Client[T])
	onDisconnect func(*Client[T], error)
	onMessage    func(T)
	onError      func(error)

	// Lifecycle
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	closeOnce sync.Once
}

// ClientOptions configures the WebSocket client
type ClientOptions struct {
	// DialOptions for the underlying connection
	DialOptions

	// Reconnection configuration
	Reconnect *ReconnectConfig

	// QueueSize is the maximum number of messages to queue during disconnection
	// 0 disables queuing
	QueueSize int

	// QueueTimeout is how long queued messages are valid
	// Default: 30 seconds
	QueueTimeout time.Duration

	// OnError is called when an error occurs
	OnError func(error)

	// OnStateChange is called when the connection state changes
	OnStateChange StateHandler
}

// DefaultClientOptions returns default client options
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		DialOptions: DialOptions{
			HandshakeTimeout: 30 * time.Second,
			ReadBufferSize:   4096,
			WriteBufferSize:  4096,
			MaxFrameSize:     4096,
			MaxMessageSize:   1048576,
			PingInterval:     30 * time.Second,
			PongTimeout:      10 * time.Second,
		},
		Reconnect:    DefaultReconnectConfig(),
		QueueSize:    100,
		QueueTimeout: 30 * time.Second,
	}
}

// NewClient creates a new WebSocket client
func NewClient[T any](url string, opts *ClientOptions) *Client[T] {
	if opts == nil {
		opts = DefaultClientOptions()
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &Client[T]{
		url:         url,
		opts:        opts,
		dialer:      NewDialer(&opts.DialOptions),
		state:       newStateManager(),
		reconnector: newReconnector(opts.Reconnect),
		ctx:         ctx,
		cancel:      cancel,
	}

	// Initialize queue if enabled
	if opts.QueueSize > 0 {
		c.queue = newMessageQueue[T](opts.QueueSize, opts.QueueTimeout)
	}

	// Register callbacks
	if opts.OnError != nil {
		c.onError = opts.OnError
	}
	if opts.OnStateChange != nil {
		c.state.OnStateChange(opts.OnStateChange)
	}

	return c
}

// OnConnect sets the callback for when the connection is established
func (c *Client[T]) OnConnect(fn func(*Client[T])) {
	c.onConnect = fn
}

// OnDisconnect sets the callback for when the connection is lost
func (c *Client[T]) OnDisconnect(fn func(*Client[T], error)) {
	c.onDisconnect = fn
}

// OnMessage sets the callback for received messages
func (c *Client[T]) OnMessage(fn func(T)) {
	c.onMessage = fn
}

// Connect establishes the WebSocket connection
func (c *Client[T]) Connect(ctx context.Context) error {
	// Transition to connecting state
	if !c.state.transition(StateDisconnected, StateConnecting, nil, 0) {
		current := c.state.State()
		if current == StateConnected {
			return nil // Already connected
		}
		return ErrInvalidState
	}

	// Attempt to connect
	conn, err := DialWithDialer[T](ctx, c.dialer, c.url)
	if err != nil {
		c.state.forceTransition(StateDisconnected, err, 0)
		return err
	}

	c.connMu.Lock()
	c.conn = conn
	c.connMu.Unlock()

	// Transition to connected state
	c.state.forceTransition(StateConnected, nil, 0)

	// Flush any queued messages
	if c.queue != nil {
		c.queue.Flush(func(ctx context.Context, msg T) error {
			return c.write(ctx, msg)
		})
	}

	// Call connect callback
	if c.onConnect != nil {
		c.onConnect(c)
	}

	// Reset reconnection attempts after successful connection
	c.reconnector.maybeReset()

	return nil
}

// ConnectWithReadLoop connects and starts a read loop
// Messages are delivered via OnMessage callback
func (c *Client[T]) ConnectWithReadLoop(ctx context.Context) error {
	if err := c.Connect(ctx); err != nil {
		return err
	}

	c.wg.Add(1)
	go c.readLoop()

	return nil
}

// readLoop continuously reads messages from the connection
func (c *Client[T]) readLoop() {
	defer c.wg.Done()

	for {
		// Check if client is closed
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		// Check connection state
		state := c.state.State()
		if state == StateClosed || state == StateClosing {
			return
		}

		// Skip if not connected
		if state != StateConnected {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Read message
		msg, err := c.Read(c.ctx)
		if err != nil {
			// Handle disconnection - check for CloseError, ErrConnectionClosed, or context canceled
			if IsCloseError(err) || err == ErrConnectionClosed || err == ErrContextCanceled {
				c.handleDisconnect(err)
				continue
			}

			// Report error
			if c.onError != nil {
				c.onError(err)
			}
			continue
		}

		// Deliver message
		if c.onMessage != nil {
			c.onMessage(msg)
		}
	}
}

// handleDisconnect handles connection loss
func (c *Client[T]) handleDisconnect(err error) {
	// Transition to disconnected or reconnecting
	fromState := c.state.State()
	if fromState == StateClosed || fromState == StateClosing {
		return
	}

	// Call disconnect callback
	if c.onDisconnect != nil {
		c.onDisconnect(c, err)
	}

	// Check if we should reconnect
	if c.reconnector.shouldReconnect(err) {
		c.state.forceTransition(StateReconnecting, err, c.reconnector.attempts)
		c.startReconnect()
	} else {
		c.state.forceTransition(StateDisconnected, err, 0)
	}
}

// startReconnect initiates the reconnection process
func (c *Client[T]) startReconnect() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		err := c.reconnector.reconnectLoop(c.ctx, func(ctx context.Context) error {
			// Transition to connecting
			c.state.forceTransition(StateConnecting, nil, c.reconnector.attempts)

			// Attempt connection
			conn, err := DialWithDialer[T](ctx, c.dialer, c.url)
			if err != nil {
				return err
			}

			c.connMu.Lock()
			c.conn = conn
			c.connMu.Unlock()

			// Transition to connected
			c.state.forceTransition(StateConnected, nil, c.reconnector.attempts)

			// Flush queued messages
			if c.queue != nil {
				c.queue.Flush(func(ctx context.Context, msg T) error {
					return c.write(ctx, msg)
				})
			}

			// Call connect callback
			if c.onConnect != nil {
				c.onConnect(c)
			}

			return nil
		})

		if err != nil {
			c.state.forceTransition(StateDisconnected, err, c.reconnector.attempts)
			if c.onError != nil {
				c.onError(err)
			}
		}
	}()
}

// Read reads a message from the connection
func (c *Client[T]) Read(ctx context.Context) (T, error) {
	var zero T

	c.connMu.RLock()
	conn := c.conn
	c.connMu.RUnlock()

	if conn == nil {
		return zero, ErrConnectionClosed
	}

	return conn.Read(ctx)
}

// Write writes a message to the connection
// If disconnected and queue is enabled, the message is queued
func (c *Client[T]) Write(ctx context.Context, msg T) error {
	state := c.state.State()

	// If connected, send immediately
	if state == StateConnected {
		return c.write(ctx, msg)
	}

	// If queue is enabled, queue the message
	if c.queue != nil && (state == StateReconnecting || state == StateConnecting) {
		errCh, err := c.queue.Enqueue(ctx, msg)
		if err != nil {
			return err
		}

		// Wait for result or context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	}

	return ErrConnectionClosed
}

// write sends a message directly to the connection
func (c *Client[T]) write(ctx context.Context, msg T) error {
	c.connMu.RLock()
	conn := c.conn
	c.connMu.RUnlock()

	if conn == nil {
		return ErrConnectionClosed
	}

	return conn.Write(ctx, msg)
}

// State returns the current connection state
func (c *Client[T]) State() ConnectionState {
	return c.state.State()
}

// IsConnected returns true if the client is connected
func (c *Client[T]) IsConnected() bool {
	return c.state.State() == StateConnected
}

// OnStateChange registers a callback for state change events
func (c *Client[T]) OnStateChange(handler StateHandler) {
	c.state.OnStateChange(handler)
}

// SetSessionID sets the session identifier for reconnection
func (c *Client[T]) SetSessionID(id string) {
	c.state.SetSessionID(id)
}

// SessionID returns the current session identifier
func (c *Client[T]) SessionID() string {
	return c.state.SessionID()
}

// QueueStats returns message queue statistics
func (c *Client[T]) QueueStats() MessageQueueStats {
	if c.queue == nil {
		return MessageQueueStats{}
	}
	return c.queue.Stats()
}

// Close closes the client and underlying connection
func (c *Client[T]) Close() error {
	var closeErr error

	c.closeOnce.Do(func() {
		// Transition to closing state
		c.state.forceTransition(StateClosing, nil, 0)

		// Cancel context to stop all goroutines
		c.cancel()

		// Close the queue
		if c.queue != nil {
			c.queue.Close()
		}

		// Close the connection
		c.connMu.Lock()
		if c.conn != nil {
			closeErr = c.conn.Close(1000, "client closed")
			c.conn = nil
		}
		c.connMu.Unlock()

		// Wait for goroutines to finish
		c.wg.Wait()

		// Transition to closed state
		c.state.forceTransition(StateClosed, nil, 0)
	})

	return closeErr
}

// Conn returns the underlying connection (may be nil if disconnected)
func (c *Client[T]) Conn() *Conn[T] {
	c.connMu.RLock()
	defer c.connMu.RUnlock()
	return c.conn
}
