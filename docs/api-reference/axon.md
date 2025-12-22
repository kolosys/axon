# axon API

Complete API documentation for the axon package.

**Import Path:** `github.com/kolosys/axon`

## Package Documentation



## Variables

**ErrInvalidHandshake, ErrUpgradeRequired, ErrInvalidOrigin, ErrInvalidSubprotocol, ErrConnectionClosed, ErrFrameTooLarge, ErrMessageTooLarge, ErrInvalidFrame, ErrInvalidMask, ErrUnsupportedFrameType, ErrFragmentedControlFrame, ErrInvalidCloseCode, ErrReadDeadlineExceeded, ErrWriteDeadlineExceeded, ErrContextCanceled, ErrSerializationFailed, ErrDeserializationFailed, ErrReconnectFailed, ErrQueueFull, ErrQueueClosed, ErrQueueTimeout, ErrQueueCleared, ErrCompressionFailed, ErrInvalidState, ErrClientClosed**

Sentinel errors for WebSocket protocol violations and state errors


```go
var ErrInvalidHandshake = errors.New("axon: invalid websocket handshake")	// ErrInvalidHandshake indicates the WebSocket handshake failed

var ErrUpgradeRequired = errors.New("axon: upgrade required")	// ErrUpgradeRequired indicates the request is not a valid WebSocket upgrade request

var ErrInvalidOrigin = errors.New("axon: invalid origin")	// ErrInvalidOrigin indicates the origin header is not allowed

var ErrInvalidSubprotocol = errors.New("axon: invalid subprotocol")	// ErrInvalidSubprotocol indicates the requested subprotocol is not supported

var ErrConnectionClosed = errors.New("axon: connection closed")	// ErrConnectionClosed indicates the connection has been closed

var ErrFrameTooLarge = errors.New("axon: frame too large")	// ErrFrameTooLarge indicates a frame exceeds the maximum allowed size

var ErrMessageTooLarge = errors.New("axon: message too large")	// ErrMessageTooLarge indicates a message exceeds the maximum allowed size

var ErrInvalidFrame = errors.New("axon: invalid frame")	// ErrInvalidFrame indicates a frame violates the WebSocket protocol

var ErrInvalidMask = errors.New("axon: invalid mask")	// ErrInvalidMask indicates frame masking is invalid

var ErrUnsupportedFrameType = errors.New("axon: unsupported frame type")	// ErrUnsupportedFrameType indicates the frame type is not supported

var ErrFragmentedControlFrame = errors.New("axon: fragmented control frame")	// ErrFragmentedControlFrame indicates a control frame is fragmented (not allowed)

var ErrInvalidCloseCode = errors.New("axon: invalid close code")	// ErrInvalidCloseCode indicates an invalid close code was used

var ErrReadDeadlineExceeded = errors.New("axon: read deadline exceeded")	// ErrReadDeadlineExceeded indicates a read operation exceeded its deadline

var ErrWriteDeadlineExceeded = errors.New("axon: write deadline exceeded")	// ErrWriteDeadlineExceeded indicates a write operation exceeded its deadline

var ErrContextCanceled = errors.New("axon: context canceled")	// ErrContextCanceled indicates the context was canceled

var ErrSerializationFailed = errors.New("axon: serialization failed")	// ErrSerializationFailed indicates message serialization failed

var ErrDeserializationFailed = errors.New("axon: deserialization failed")	// ErrDeserializationFailed indicates message deserialization failed

var ErrReconnectFailed = errors.New("axon: reconnection failed")	// ErrReconnectFailed indicates reconnection attempts have been exhausted

var ErrQueueFull = errors.New("axon: message queue full")	// ErrQueueFull indicates the message queue is full

var ErrQueueClosed = errors.New("axon: message queue closed")	// ErrQueueClosed indicates the message queue has been closed

var ErrQueueTimeout = errors.New("axon: queued message timeout")	// ErrQueueTimeout indicates a queued message has expired

var ErrQueueCleared = errors.New("axon: message queue cleared")	// ErrQueueCleared indicates the queue was cleared before message was sent

var ErrCompressionFailed = errors.New("axon: compression failed")	// ErrCompressionFailed indicates compression or decompression failed

var ErrInvalidState = errors.New("axon: invalid state transition")	// ErrInvalidState indicates an invalid state transition was attempted

var ErrClientClosed = errors.New("axon: client closed")	// ErrClientClosed indicates the client has been closed

```

**DefaultMetrics**

DefaultMetrics is the default metrics instance


```go
var DefaultMetrics = &Metrics{}
```

## Types

### Client
Client is a WebSocket client with automatic reconnection and message queuing

#### Example Usage

```go
// Create a new Client
client := Client{

}
```

#### Type Definition

```go
type Client struct {
}
```

### Constructor Functions

### NewClient

NewClient creates a new WebSocket client

```go
func NewClient(url string, opts *ClientOptions) **ast.IndexExpr
```

**Parameters:**
- `url` (string)
- `opts` (*ClientOptions)

**Returns:**
- **ast.IndexExpr

## Methods

### Close

Close closes the client and underlying connection

```go
func (*CompressionManager) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Conn

Conn returns the underlying connection (may be nil if disconnected)

```go
func (**ast.IndexExpr) Conn() **ast.IndexExpr
```

**Parameters:**
  None

**Returns:**
- **ast.IndexExpr

### Connect

Connect establishes the WebSocket connection

```go
func (**ast.IndexExpr) Connect(ctx context.Context) error
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- error

### ConnectWithReadLoop

ConnectWithReadLoop connects and starts a read loop Messages are delivered via OnMessage callback

```go
func (**ast.IndexExpr) ConnectWithReadLoop(ctx context.Context) error
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- error

### IsConnected

IsConnected returns true if the client is connected

```go
func (**ast.IndexExpr) IsConnected() bool
```

**Parameters:**
  None

**Returns:**
- bool

### OnConnect

OnConnect sets the callback for when the connection is established

```go
func (**ast.IndexExpr) OnConnect(fn func(**ast.IndexExpr))
```

**Parameters:**
- `fn` (func(**ast.IndexExpr))

**Returns:**
  None

### OnDisconnect

OnDisconnect sets the callback for when the connection is lost

```go
func (**ast.IndexExpr) OnDisconnect(fn func(**ast.IndexExpr, error))
```

**Parameters:**
- `fn` (func(**ast.IndexExpr, error))

**Returns:**
  None

### OnMessage

OnMessage sets the callback for received messages

```go
func (**ast.IndexExpr) OnMessage(fn func(T))
```

**Parameters:**
- `fn` (func(T))

**Returns:**
  None

### OnStateChange

OnStateChange registers a callback for state change events

```go
func (**ast.IndexExpr) OnStateChange(handler StateHandler)
```

**Parameters:**
- `handler` (StateHandler)

**Returns:**
  None

### QueueStats

QueueStats returns message queue statistics

```go
func (**ast.IndexExpr) QueueStats() MessageQueueStats
```

**Parameters:**
  None

**Returns:**
- MessageQueueStats

### Read

Read reads a message from the connection

```go
func (**ast.IndexExpr) Read(ctx context.Context) (T, error)
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- T
- error

### SessionID

SessionID returns the current session identifier

```go
func (**ast.IndexExpr) SessionID() string
```

**Parameters:**
  None

**Returns:**
- string

### SetSessionID

SetSessionID sets the session identifier for reconnection

```go
func (**ast.IndexExpr) SetSessionID(id string)
```

**Parameters:**
- `id` (string)

**Returns:**
  None

### State

State returns the current connection state

```go
func (**ast.IndexExpr) State() ConnectionState
```

**Parameters:**
  None

**Returns:**
- ConnectionState

### Write

Write writes a message to the connection If disconnected and queue is enabled, the message is queued

```go
func (**ast.IndexExpr) Write(ctx context.Context, msg T) error
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (T)

**Returns:**
- error

### ClientOptions
ClientOptions configures the WebSocket client

#### Example Usage

```go
// Create a new ClientOptions
clientoptions := ClientOptions{
    Reconnect: &ReconnectConfig{}{},
    QueueSize: 42,
    QueueTimeout: /* value */,
    OnError: /* value */,
    OnStateChange: StateHandler{},
}
```

#### Type Definition

```go
type ClientOptions struct {
    DialOptions
    Reconnect *ReconnectConfig
    QueueSize int
    QueueTimeout time.Duration
    OnError func(error)
    OnStateChange StateHandler
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| *DialOptions | `DialOptions` | DialOptions for the underlying connection |
| Reconnect | `*ReconnectConfig` | Reconnection configuration |
| QueueSize | `int` | QueueSize is the maximum number of messages to queue during disconnection 0 disables queuing |
| QueueTimeout | `time.Duration` | QueueTimeout is how long queued messages are valid Default: 30 seconds |
| OnError | `func(error)` | OnError is called when an error occurs |
| OnStateChange | `StateHandler` | OnStateChange is called when the connection state changes |

### Constructor Functions

### DefaultClientOptions

DefaultClientOptions returns default client options

```go
func DefaultClientOptions() *ClientOptions
```

**Parameters:**
  None

**Returns:**
- *ClientOptions

### CompressionManager
CompressionManager handles per-message compression (RFC 7692)

#### Example Usage

```go
// Create a new CompressionManager
compressionmanager := CompressionManager{

}
```

#### Type Definition

```go
type CompressionManager struct {
}
```

## Methods

### Close

Close releases compression resources

```go
func (**ast.IndexExpr) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Compress

Compress compresses the payload using DEFLATE

```go
func (*CompressionManager) Compress(data []byte) ([]byte, error)
```

**Parameters:**
- `data` ([]byte)

**Returns:**
- []byte
- error

### Decompress

Decompress decompresses the payload using DEFLATE

```go
func (*CompressionManager) Decompress(data []byte) ([]byte, error)
```

**Parameters:**
- `data` ([]byte)

**Returns:**
- []byte
- error

### ShouldCompress

ShouldCompress returns true if the payload should be compressed

```go
func (*CompressionManager) ShouldCompress(payloadSize int) bool
```

**Parameters:**
- `payloadSize` (int)

**Returns:**
- bool

### Conn
Conn represents a WebSocket connection with type-safe message handling

#### Example Usage

```go
// Create a new Conn
conn := Conn{

}
```

#### Type Definition

```go
type Conn struct {
}
```

### Constructor Functions

### Dial

Dial connects to a WebSocket server at the given URL. The URL must have a ws:// or wss:// scheme.

```go
func Dial(ctx context.Context, rawURL string, opts *DialOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `ctx` (context.Context)
- `rawURL` (string)
- `opts` (*DialOptions)

**Returns:**
- **ast.IndexExpr
- error

### DialWithDialer

DialWithDialer connects to a WebSocket server using the provided dialer

```go
func DialWithDialer(ctx context.Context, d *Dialer, rawURL string) (**ast.IndexExpr, error)
```

**Parameters:**
- `ctx` (context.Context)
- `d` (*Dialer)
- `rawURL` (string)

**Returns:**
- **ast.IndexExpr
- error

### Upgrade

Upgrade upgrades an HTTP connection to a WebSocket connection

```go
func Upgrade(w http.ResponseWriter, r *http.Request, opts *UpgradeOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `w` (http.ResponseWriter)
- `r` (*http.Request)
- `opts` (*UpgradeOptions)

**Returns:**
- **ast.IndexExpr
- error

## Methods

### Close

Close closes the connection with the given code and reason

```go
func (**ast.IndexExpr) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### CloseCode

CloseCode returns the close code if the connection was closed

```go
func (**ast.IndexExpr) CloseCode() int
```

**Parameters:**
  None

**Returns:**
- int

### CloseReason

CloseReason returns the close reason if the connection was closed

```go
func (**ast.IndexExpr) CloseReason() string
```

**Parameters:**
  None

**Returns:**
- string

### IsClosed

IsClosed returns true if the connection has been closed

```go
func (**ast.IndexExpr) IsClosed() bool
```

**Parameters:**
  None

**Returns:**
- bool

### Read

Read reads a complete message from the connection

```go
func (**ast.IndexExpr) Read(ctx context.Context) (T, error)
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- T
- error

### Write

Write writes a message to the connection

```go
func (**ast.IndexExpr) Write(ctx context.Context, msg T) error
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (T)

**Returns:**
- error

### ConnectionState
ConnectionState represents the state of a WebSocket connection

#### Example Usage

```go
// Example usage of ConnectionState
var value ConnectionState
// Initialize with appropriate value
```

#### Type Definition

```go
type ConnectionState int32
```

## Methods

### CanReconnect

CanReconnect returns true if reconnection is allowed from this state

```go
func (ConnectionState) CanReconnect() bool
```

**Parameters:**
  None

**Returns:**
- bool

### IsActive

IsActive returns true if the state represents an active or transitioning state

```go
func (ConnectionState) IsActive() bool
```

**Parameters:**
  None

**Returns:**
- bool

### String

String returns the string representation of the connection state

```go
func (ConnectionState) String() string
```

**Parameters:**
  None

**Returns:**
- string

### DialOptions
DialOptions configures WebSocket client dial options

#### Example Usage

```go
// Create a new DialOptions
dialoptions := DialOptions{
    HandshakeTimeout: /* value */,
    ReadBufferSize: 42,
    WriteBufferSize: 42,
    MaxFrameSize: 42,
    MaxMessageSize: 42,
    ReadDeadline: /* value */,
    WriteDeadline: /* value */,
    PingInterval: /* value */,
    PongTimeout: /* value */,
    Subprotocols: [],
    Compression: true,
    CompressionThreshold: 42,
    Headers: /* value */,
    TLSConfig: &/* value */{},
    NetDialer: &/* value */{},
}
```

#### Type Definition

```go
type DialOptions struct {
    HandshakeTimeout time.Duration
    ReadBufferSize int
    WriteBufferSize int
    MaxFrameSize int
    MaxMessageSize int
    ReadDeadline time.Duration
    WriteDeadline time.Duration
    PingInterval time.Duration
    PongTimeout time.Duration
    Subprotocols []string
    Compression bool
    CompressionThreshold int
    Headers http.Header
    TLSConfig *tls.Config
    NetDialer *net.Dialer
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| HandshakeTimeout | `time.Duration` | HandshakeTimeout is the maximum time to wait for the handshake to complete. Default is 30 seconds. |
| ReadBufferSize | `int` | ReadBufferSize sets the size of the read buffer in bytes. Default is 4096 bytes. |
| WriteBufferSize | `int` | WriteBufferSize sets the size of the write buffer in bytes. Default is 4096 bytes. |
| MaxFrameSize | `int` | MaxFrameSize sets the maximum frame size in bytes. Default is 4096 bytes. |
| MaxMessageSize | `int` | MaxMessageSize sets the maximum message size in bytes. Default is 1048576 bytes (1MB). |
| ReadDeadline | `time.Duration` | ReadDeadline sets the read deadline for connections. Default is no deadline. |
| WriteDeadline | `time.Duration` | WriteDeadline sets the write deadline for connections. Default is no deadline. |
| PingInterval | `time.Duration` | PingInterval sets the interval for sending ping frames. If zero, pings are disabled. |
| PongTimeout | `time.Duration` | PongTimeout sets the timeout for waiting for a pong response. If zero, pong timeout is disabled. |
| Subprotocols | `[]string` | Subprotocols sets the list of supported subprotocols. Default is nil (no subprotocols). |
| Compression | `bool` | Compression enables per-message compression (RFC 7692). Default is false (disabled). |
| CompressionThreshold | `int` | CompressionThreshold sets the minimum message size to compress. Messages smaller than this will not be compressed. Default is 256 bytes. |
| Headers | `http.Header` | Headers sets additional HTTP headers for the handshake request. |
| TLSConfig | `*tls.Config` | TLSConfig specifies the TLS configuration to use for wss:// connections. If nil, the default configuration is used. |
| NetDialer | `*net.Dialer` | NetDialer specifies the dialer to use for creating the network connection. If nil, a default dialer is used. |

### Dialer
Dialer is a WebSocket client dialer

#### Example Usage

```go
// Create a new Dialer
dialer := Dialer{

}
```

#### Type Definition

```go
type Dialer struct {
}
```

### Constructor Functions

### NewDialer

NewDialer creates a new Dialer with the given options

```go
func NewDialer(opts *DialOptions) *Dialer
```

**Parameters:**
- `opts` (*DialOptions)

**Returns:**
- *Dialer

### Frame
Frame represents a WebSocket frame

#### Example Usage

```go
// Create a new Frame
frame := Frame{
    Fin: true,
    Rsv1: true,
    Rsv2: true,
    Rsv3: true,
    Opcode: byte{},
    Masked: true,
    MaskKey: [],
    Payload: [],
}
```

#### Type Definition

```go
type Frame struct {
    Fin bool
    Rsv1 bool
    Rsv2 bool
    Rsv3 bool
    Opcode byte
    Masked bool
    MaskKey []byte
    Payload []byte
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Fin | `bool` |  |
| Rsv1 | `bool` |  |
| Rsv2 | `bool` |  |
| Rsv3 | `bool` |  |
| Opcode | `byte` |  |
| Masked | `bool` |  |
| MaskKey | `[]byte` |  |
| Payload | `[]byte` |  |

### MessageQueue
MessageQueue manages queuing of messages during disconnection

#### Example Usage

```go
// Create a new MessageQueue
messagequeue := MessageQueue{

}
```

#### Type Definition

```go
type MessageQueue struct {
}
```

## Methods

### Clear

Clear discards all queued messages

```go
func (**ast.IndexExpr) Clear()
```

**Parameters:**
  None

**Returns:**
  None

### Close

Close closes the queue and discards all pending messages

```go
func (**ast.IndexExpr) Close()
```

**Parameters:**
  None

**Returns:**
  None

### Enqueue

Enqueue adds a message to the queue Returns an error channel that will receive the send result

```go
func (**ast.IndexExpr) Enqueue(ctx context.Context, msg T) (chan error, error)
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (T)

**Returns:**
- chan error
- error

### Flush

Flush sends all queued messages using the provided send function

```go
func (**ast.IndexExpr) Flush(sendFn func(context.Context, T) error)
```

**Parameters:**
- `sendFn` (func(context.Context, T) error)

**Returns:**
  None

### Size

Size returns the current number of queued messages

```go
func (**ast.IndexExpr) Size() int
```

**Parameters:**
  None

**Returns:**
- int

### Stats

Stats returns queue statistics

```go
func (**ast.IndexExpr) Stats() MessageQueueStats
```

**Parameters:**
  None

**Returns:**
- MessageQueueStats

### MessageQueueStats
MessageQueueStats contains queue statistics

#### Example Usage

```go
// Create a new MessageQueueStats
messagequeuestats := MessageQueueStats{
    CurrentSize: 42,
    MaxSize: 42,
    Dropped: 42,
    Enqueued: 42,
    Sent: 42,
}
```

#### Type Definition

```go
type MessageQueueStats struct {
    CurrentSize int
    MaxSize int
    Dropped int64
    Enqueued int64
    Sent int64
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| CurrentSize | `int` |  |
| MaxSize | `int` |  |
| Dropped | `int64` |  |
| Enqueued | `int64` |  |
| Sent | `int64` |  |

### Metrics
Metrics tracks WebSocket connection metrics

#### Example Usage

```go
// Create a new Metrics
metrics := Metrics{
    ActiveConnections: /* value */,
    TotalConnections: /* value */,
    ClosedConnections: /* value */,
    MessagesRead: /* value */,
    MessagesWritten: /* value */,
    BytesRead: /* value */,
    BytesWritten: /* value */,
    ReadErrors: /* value */,
    WriteErrors: /* value */,
    FrameErrors: /* value */,
    HandshakeErrors: /* value */,
    ReadLatency: /* value */,
    WriteLatency: /* value */,
    ReconnectAttempts: /* value */,
    ReconnectSuccesses: /* value */,
    ReconnectFailures: /* value */,
    QueueEnqueued: /* value */,
    QueueSent: /* value */,
    QueueDropped: /* value */,
    CompressedMessages: /* value */,
    DecompressedMessages: /* value */,
    CompressionSaved: /* value */,
}
```

#### Type Definition

```go
type Metrics struct {
    ActiveConnections atomic.Int64
    TotalConnections atomic.Int64
    ClosedConnections atomic.Int64
    MessagesRead atomic.Int64
    MessagesWritten atomic.Int64
    BytesRead atomic.Int64
    BytesWritten atomic.Int64
    ReadErrors atomic.Int64
    WriteErrors atomic.Int64
    FrameErrors atomic.Int64
    HandshakeErrors atomic.Int64
    ReadLatency atomic.Int64
    WriteLatency atomic.Int64
    ReconnectAttempts atomic.Int64
    ReconnectSuccesses atomic.Int64
    ReconnectFailures atomic.Int64
    QueueEnqueued atomic.Int64
    QueueSent atomic.Int64
    QueueDropped atomic.Int64
    CompressedMessages atomic.Int64
    DecompressedMessages atomic.Int64
    CompressionSaved atomic.Int64
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ActiveConnections | `atomic.Int64` | Connection metrics |
| TotalConnections | `atomic.Int64` |  |
| ClosedConnections | `atomic.Int64` |  |
| MessagesRead | `atomic.Int64` | Message metrics |
| MessagesWritten | `atomic.Int64` |  |
| BytesRead | `atomic.Int64` |  |
| BytesWritten | `atomic.Int64` |  |
| ReadErrors | `atomic.Int64` | Error metrics |
| WriteErrors | `atomic.Int64` |  |
| FrameErrors | `atomic.Int64` |  |
| HandshakeErrors | `atomic.Int64` |  |
| ReadLatency | `atomic.Int64` | Performance metrics |
| WriteLatency | `atomic.Int64` | nanoseconds |
| ReconnectAttempts | `atomic.Int64` | Reconnection metrics |
| ReconnectSuccesses | `atomic.Int64` |  |
| ReconnectFailures | `atomic.Int64` |  |
| QueueEnqueued | `atomic.Int64` | Queue metrics |
| QueueSent | `atomic.Int64` |  |
| QueueDropped | `atomic.Int64` |  |
| CompressedMessages | `atomic.Int64` | Compression metrics |
| DecompressedMessages | `atomic.Int64` |  |
| CompressionSaved | `atomic.Int64` | bytes saved by compression |

## Methods

### GetSnapshot

GetSnapshot returns a snapshot of current metrics

```go
func (*Metrics) GetSnapshot() MetricsSnapshot
```

**Parameters:**
  None

**Returns:**
- MetricsSnapshot

### RecordCompression

RecordCompression records a compression operation

```go
func (*Metrics) RecordCompression(originalSize, compressedSize int)
```

**Parameters:**
- `originalSize` (int)
- `compressedSize` (int)

**Returns:**
  None

### RecordConnection

RecordConnection records a new connection

```go
func (*Metrics) RecordConnection()
```

**Parameters:**
  None

**Returns:**
  None

### RecordDecompression

RecordDecompression records a decompression operation

```go
func (*Metrics) RecordDecompression()
```

**Parameters:**
  None

**Returns:**
  None

### RecordDisconnection

RecordDisconnection records a connection closure

```go
func (*Metrics) RecordDisconnection()
```

**Parameters:**
  None

**Returns:**
  None

### RecordFrameError

RecordFrameError records a frame parsing error

```go
func (*Metrics) RecordFrameError()
```

**Parameters:**
  None

**Returns:**
  None

### RecordHandshakeError

RecordHandshakeError records a handshake error

```go
func (*Metrics) RecordHandshakeError()
```

**Parameters:**
  None

**Returns:**
  None

### RecordQueueDropped

RecordQueueDropped records a queued message being dropped

```go
func (*Metrics) RecordQueueDropped()
```

**Parameters:**
  None

**Returns:**
  None

### RecordQueueEnqueue

RecordQueueEnqueue records a message being queued

```go
func (*Metrics) RecordQueueEnqueue()
```

**Parameters:**
  None

**Returns:**
  None

### RecordQueueSent

RecordQueueSent records a queued message being sent

```go
func (*Metrics) RecordQueueSent()
```

**Parameters:**
  None

**Returns:**
  None

### RecordRead

RecordRead records a read operation

```go
func (*Metrics) RecordRead(bytes int, latency time.Duration)
```

**Parameters:**
- `bytes` (int)
- `latency` (time.Duration)

**Returns:**
  None

### RecordReadError

RecordReadError records a read error

```go
func (*Metrics) RecordReadError()
```

**Parameters:**
  None

**Returns:**
  None

### RecordReconnectAttempt

RecordReconnectAttempt records a reconnection attempt

```go
func (*Metrics) RecordReconnectAttempt()
```

**Parameters:**
  None

**Returns:**
  None

### RecordReconnectFailure

RecordReconnectFailure records a failed reconnection

```go
func (*Metrics) RecordReconnectFailure()
```

**Parameters:**
  None

**Returns:**
  None

### RecordReconnectSuccess

RecordReconnectSuccess records a successful reconnection

```go
func (*Metrics) RecordReconnectSuccess()
```

**Parameters:**
  None

**Returns:**
  None

### RecordWrite

RecordWrite records a write operation

```go
func (*Metrics) RecordWrite(bytes int, latency time.Duration)
```

**Parameters:**
- `bytes` (int)
- `latency` (time.Duration)

**Returns:**
  None

### RecordWriteError

RecordWriteError records a write error

```go
func (*Metrics) RecordWriteError()
```

**Parameters:**
  None

**Returns:**
  None

### MetricsSnapshot
MetricsSnapshot represents a snapshot of metrics at a point in time

#### Example Usage

```go
// Create a new MetricsSnapshot
metricssnapshot := MetricsSnapshot{
    ActiveConnections: 42,
    TotalConnections: 42,
    ClosedConnections: 42,
    MessagesRead: 42,
    MessagesWritten: 42,
    BytesRead: 42,
    BytesWritten: 42,
    ReadErrors: 42,
    WriteErrors: 42,
    FrameErrors: 42,
    HandshakeErrors: 42,
    AvgReadLatency: /* value */,
    AvgWriteLatency: /* value */,
    ReconnectAttempts: 42,
    ReconnectSuccesses: 42,
    ReconnectFailures: 42,
    QueueEnqueued: 42,
    QueueSent: 42,
    QueueDropped: 42,
    CompressedMessages: 42,
    DecompressedMessages: 42,
    CompressionSaved: 42,
}
```

#### Type Definition

```go
type MetricsSnapshot struct {
    ActiveConnections int64
    TotalConnections int64
    ClosedConnections int64
    MessagesRead int64
    MessagesWritten int64
    BytesRead int64
    BytesWritten int64
    ReadErrors int64
    WriteErrors int64
    FrameErrors int64
    HandshakeErrors int64
    AvgReadLatency time.Duration
    AvgWriteLatency time.Duration
    ReconnectAttempts int64
    ReconnectSuccesses int64
    ReconnectFailures int64
    QueueEnqueued int64
    QueueSent int64
    QueueDropped int64
    CompressedMessages int64
    DecompressedMessages int64
    CompressionSaved int64
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ActiveConnections | `int64` |  |
| TotalConnections | `int64` |  |
| ClosedConnections | `int64` |  |
| MessagesRead | `int64` |  |
| MessagesWritten | `int64` |  |
| BytesRead | `int64` |  |
| BytesWritten | `int64` |  |
| ReadErrors | `int64` |  |
| WriteErrors | `int64` |  |
| FrameErrors | `int64` |  |
| HandshakeErrors | `int64` |  |
| AvgReadLatency | `time.Duration` |  |
| AvgWriteLatency | `time.Duration` |  |
| ReconnectAttempts | `int64` | Reconnection metrics |
| ReconnectSuccesses | `int64` |  |
| ReconnectFailures | `int64` |  |
| QueueEnqueued | `int64` | Queue metrics |
| QueueSent | `int64` |  |
| QueueDropped | `int64` |  |
| CompressedMessages | `int64` | Compression metrics |
| DecompressedMessages | `int64` |  |
| CompressionSaved | `int64` |  |

### ReconnectConfig
ReconnectConfig configures automatic reconnection behavior

#### Example Usage

```go
// Create a new ReconnectConfig
reconnectconfig := ReconnectConfig{
    Enabled: true,
    MaxAttempts: 42,
    InitialDelay: /* value */,
    MaxDelay: /* value */,
    BackoffMultiplier: 3.14,
    Jitter: true,
    ResetAfter: /* value */,
    OnReconnecting: /* value */,
    OnReconnected: /* value */,
    OnReconnectFailed: /* value */,
    ShouldReconnect: /* value */,
}
```

#### Type Definition

```go
type ReconnectConfig struct {
    Enabled bool
    MaxAttempts int
    InitialDelay time.Duration
    MaxDelay time.Duration
    BackoffMultiplier float64
    Jitter bool
    ResetAfter time.Duration
    OnReconnecting func(attempt int, delay time.Duration)
    OnReconnected func(attempt int)
    OnReconnectFailed func(attempt int, err error)
    ShouldReconnect func(err error, attempt int) bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Enabled | `bool` | Enabled determines if automatic reconnection is active |
| MaxAttempts | `int` | MaxAttempts is the maximum number of reconnection attempts (0 = unlimited) |
| InitialDelay | `time.Duration` | InitialDelay is the initial delay before the first reconnection attempt Default: 1 second |
| MaxDelay | `time.Duration` | MaxDelay is the maximum delay between reconnection attempts Default: 30 seconds |
| BackoffMultiplier | `float64` | BackoffMultiplier is the multiplier for exponential backoff Default: 2.0 |
| Jitter | `bool` | Jitter adds randomness to delay to prevent thundering herd Default: true |
| ResetAfter | `time.Duration` | ResetAfter resets the attempt counter after being connected for this duration Default: 60 seconds |
| OnReconnecting | `func(attempt int, delay time.Duration)` | OnReconnecting is called when a reconnection attempt starts |
| OnReconnected | `func(attempt int)` | OnReconnected is called when reconnection succeeds |
| OnReconnectFailed | `func(attempt int, err error)` | OnReconnectFailed is called when reconnection fails |
| ShouldReconnect | `func(err error, attempt int) bool` | ShouldReconnect is called to determine if reconnection should be attempted If nil, always attempts to reconnect (within MaxAttempts) |

### Constructor Functions

### DefaultReconnectConfig

DefaultReconnectConfig returns a default reconnection configuration

```go
func DefaultReconnectConfig() *ReconnectConfig
```

**Parameters:**
  None

**Returns:**
- *ReconnectConfig

### StateChange
StateChange represents a state transition event

#### Example Usage

```go
// Create a new StateChange
statechange := StateChange{
    From: ConnectionState{},
    To: ConnectionState{},
    Time: /* value */,
    Err: error{},
    Attempt: 42,
    SessionID: "example",
}
```

#### Type Definition

```go
type StateChange struct {
    From ConnectionState
    To ConnectionState
    Time time.Time
    Err error
    Attempt int
    SessionID string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| From | `ConnectionState` |  |
| To | `ConnectionState` |  |
| Time | `time.Time` |  |
| Err | `error` | Error that caused the transition (if any) |
| Attempt | `int` | Reconnection attempt number (if applicable) |
| SessionID | `string` | Current session identifier |

### StateHandler
StateHandler is a callback for state change events

#### Example Usage

```go
// Example usage of StateHandler
var value StateHandler
// Initialize with appropriate value
```

#### Type Definition

```go
type StateHandler func(StateChange)
```

### UpgradeOptions
UpgradeOptions configures WebSocket connection upgrade options

#### Example Usage

```go
// Create a new UpgradeOptions
upgradeoptions := UpgradeOptions{
    ReadBufferSize: 42,
    WriteBufferSize: 42,
    MaxFrameSize: 42,
    MaxMessageSize: 42,
    ReadDeadline: /* value */,
    WriteDeadline: /* value */,
    PingInterval: /* value */,
    PongTimeout: /* value */,
    CheckOrigin: /* value */,
    Subprotocols: [],
    Compression: true,
}
```

#### Type Definition

```go
type UpgradeOptions struct {
    ReadBufferSize int
    WriteBufferSize int
    MaxFrameSize int
    MaxMessageSize int
    ReadDeadline time.Duration
    WriteDeadline time.Duration
    PingInterval time.Duration
    PongTimeout time.Duration
    CheckOrigin func(r *http.Request) bool
    Subprotocols []string
    Compression bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ReadBufferSize | `int` | ReadBufferSize sets the size of the read buffer in bytes. Default is 4096 bytes. |
| WriteBufferSize | `int` | WriteBufferSize sets the size of the write buffer in bytes. Default is 4096 bytes. |
| MaxFrameSize | `int` | MaxFrameSize sets the maximum frame size in bytes. Frames exceeding this size will result in ErrFrameTooLarge. Default is 4096 bytes. |
| MaxMessageSize | `int` | MaxMessageSize sets the maximum message size in bytes. Messages exceeding this size will result in ErrMessageTooLarge. Default is 1048576 bytes (1MB). |
| ReadDeadline | `time.Duration` | ReadDeadline sets the read deadline for connections. Default is no deadline. |
| WriteDeadline | `time.Duration` | WriteDeadline sets the write deadline for connections. Default is no deadline. |
| PingInterval | `time.Duration` | PingInterval sets the interval for sending ping frames. If zero, pings are disabled. Default is 0 (disabled). |
| PongTimeout | `time.Duration` | PongTimeout sets the timeout for waiting for a pong response. If zero, pong timeout is disabled. Default is 0 (disabled). |
| CheckOrigin | `func(r *http.Request) bool` | CheckOrigin sets a function to validate the origin header. If nil, all origins are allowed. Default is nil (all origins allowed). |
| Subprotocols | `[]string` | Subprotocols sets the list of supported subprotocols. The client's requested subprotocol must match one of these. Default is nil (no subprotocols). |
| Compression | `bool` | Compression enables per-message compression (RFC 7692). Default is false (disabled). |

### Upgrader
Upgrader handles WebSocket connection upgrades

#### Example Usage

```go
// Create a new Upgrader
upgrader := Upgrader{

}
```

#### Type Definition

```go
type Upgrader struct {
}
```

### Constructor Functions

### NewUpgrader

NewUpgrader creates a new Upgrader with default settings

```go
func NewUpgrader(opts *UpgradeOptions) *Upgrader
```

**Parameters:**
- `opts` (*UpgradeOptions)

**Returns:**
- *Upgrader

## External Links

- [Package Overview](../packages/axon.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/axon)
- [Source Code](https://github.com/kolosys/axon/tree/main/github.com/kolosys/axon)
