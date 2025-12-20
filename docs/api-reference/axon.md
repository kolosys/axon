# axon API

Complete API documentation for the axon package.

**Import Path:** `github.com/kolosys/axon`

## Package Documentation



## Variables

**ErrInvalidHandshake, ErrUpgradeRequired, ErrInvalidOrigin, ErrInvalidSubprotocol, ErrConnectionClosed, ErrFrameTooLarge, ErrMessageTooLarge, ErrInvalidFrame, ErrInvalidMask, ErrUnsupportedFrameType, ErrFragmentedControlFrame, ErrInvalidCloseCode, ErrReadDeadlineExceeded, ErrWriteDeadlineExceeded, ErrContextCanceled, ErrSerializationFailed, ErrDeserializationFailed**

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

```

**DefaultMetrics**

DefaultMetrics is the default metrics instance


```go
var DefaultMetrics = &Metrics{}
```

## Types

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
func (**ast.IndexExpr) Close(code int, reason string) error
```

**Parameters:**
- `code` (int)
- `reason` (string)

**Returns:**
- error

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

### RecordConnection

RecordConnection records a new connection

```go
func (*Metrics) RecordConnection()
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
