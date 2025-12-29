# mock API

Complete API documentation for the mock package.

**Import Path:** `github.com/kolosys/axon/mock`

## Package Documentation

Package mock provides in-memory mocks for testing code that uses the axon WebSocket library.

# Overview

The mock package offers two primary mock types:
  - MockConn[T]: Mocks the low-level Conn[T] connection for unit testing
  - MockClient[T]: Mocks the high-level Client[T] with full state machine support

Both types use in-memory channels for message passing and provide comprehensive
testing utilities including error injection, message recording, and state assertions.

# Thread Safety

All mock types are safe for concurrent use. Message recording, state management,
and error injection all use appropriate synchronization primitives (atomics,
RWMutex) for thread-safe operation.

# Usage Example - MockConn[T]

MockConn[T] is useful for testing code that directly handles WebSocket connections:

	func TestWebSocketHandler(t *testing.T) {
	    mc := mock.NewMockConn[string](nil)
	    defer mc.Close(1000, "test done")

	    // Inject a message to be read
	    mc.InjectMessage("hello")

	    // Code under test reads from connection
	    msg, err := mc.Read(context.Background())
	    if err != nil {
	        t.Fatal(err)
	    }

	    // Assert the message
	    if msg != "hello" {
	        t.Errorf("got %q, want %q", msg, "hello")
	    }

	    // Code under test writes response
	    mc.Write(context.Background(), "world")

	    // Assert what was written
	    mock.AssertWritten(t, mc, "world")
	}

# Usage Example - MockClient[T]

MockClient[T] is useful for testing high-level client behavior including
state transitions and reconnection logic:

	func TestClientReconnection(t *testing.T) {
	    mc := mock.NewMockClient[string](nil)
	    defer mc.Close()

	    // Track callbacks
	    connectCalled := mock.NewCallbackRecorder()
	    mc.OnConnect(connectCalled.Record)

	    // Connect the client
	    err := mc.Connect(context.Background())
	    if err != nil {
	        t.Fatal(err)
	    }

	    // Assert state and callback
	    mock.AssertConnected(t, mc)
	    connectCalled.AssertCalled(t)

	    // Simulate disconnection
	    mc.SimulateDisconnect(axon.ErrConnectionClosed)
	    mock.AssertDisconnected(t, mc)
	}

# Error Injection

Both mock types support error injection for testing error handling paths:

	mc := mock.NewMockConn[string](nil)
	mc.InjectReadError(axon.ErrConnectionClosed)

	_, err := mc.Read(context.Background())
	// err will be axon.ErrConnectionClosed

Clear errors when done with them:

	mc.ClearInjectedErrors()

# Message Recording

Messages are recorded by default and can be retrieved for assertions:

	mc := mock.NewMockClient[string](nil)
	mc.Connect(context.Background())

	mc.Write(context.Background(), "msg1")
	mc.Write(context.Background(), "msg2")

	msgs := mc.WrittenMessages()
	// msgs == []string{"msg1", "msg2"}

Disable recording for performance-sensitive tests:

	opts := &mock.MockClientOptions{
	    RecordMessages: false,
	}
	mc := mock.NewMockClient[string](opts)

# State Management

MockClient[T] provides full state machine support matching the real Client:

	mc := mock.NewMockClient[string](nil)

	// Track state changes
	mc.OnStateChange(func(change axon.StateChange) {
	    fmt.Printf("State: %v -> %v\n", change.From, change.To)
	})

	mc.Connect(context.Background())
	// Output: State: disconnected -> connecting
	//         State: connecting -> connected

# Assertion Helpers

The package provides helper functions for common assertions:

	// Assert connection state
	mock.AssertConnected(t, mc)
	mock.AssertDisconnected(t, mc)
	mock.AssertState(t, mc, axon.StateConnected)

	// Assert state transitions
	mock.AssertStateTransition(t, mc, axon.StateDisconnected, axon.StateConnecting)

	// Assert messages
	mock.AssertClientWritten(t, mc, expectedMsg)
	mock.AssertClientWrittenCount(t, mc, 5)

	// Wait for state or message with timeout
	mock.WaitForState(t, mc, axon.StateConnected, 5*time.Second)
	mock.WaitForClientMessage(t, mc, "expected", 1*time.Second)

# Callback Testing

Use CallbackRecorder to verify callbacks are invoked:

	mc := mock.NewMockClient[string](nil)

	connectCalled := mock.NewCallbackRecorder()
	mc.OnConnect(connectCalled.Record)

	mc.Connect(context.Background())
	connectCalled.AssertCalled(t)
	connectCalled.AssertCallCount(t, 1)

# Latency Simulation

Simulate network latency for performance testing:

	opts := &mock.MockConnOptions{
	    ReadDelay:  50 * time.Millisecond,
	    WriteDelay: 100 * time.Millisecond,
	}
	mc := mock.NewMockConn[string](opts)

# Controlled Disconnections

Simulate various disconnection scenarios:

	mc := mock.NewMockClient[string](nil)
	mc.Connect(context.Background())

	// Graceful disconnection
	mc.SimulateDisconnect(nil)

	// Disconnection with error
	mc.SimulateDisconnect(axon.ErrConnectionClosed)

	// Reconnection attempt
	mc.SimulateReconnecting(1) // attempt number 1

# Type Safety

All mocks are fully generic and maintain type safety:

	type Message struct {
	    ID   int
	    Text string
	}

	mc := mock.NewMockClient[Message](nil)
	mc.Connect(context.Background())

	msg := Message{ID: 1, Text: "hello"}
	mc.Write(context.Background(), msg)
	// Type-safe, no casting needed

# Best Practices

  - Always defer Close() on mocks to ensure cleanup
  - Use CallbackRecorder to verify callback behavior
  - Inject errors to test error handling paths
  - Use assertion helpers for readable test code
  - Disable recording in performance-critical tests
  - Clear injected errors between test cases

# Common Patterns

Testing message flow:

	mc := mock.NewMockConn[string](nil)
	defer mc.Close(1000, "")

	mc.InjectMessage("server message")
	msg, _ := mc.Read(context.Background())
	mock.AssertWritten(t, mc, expectedReply)

Testing state transitions:

	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.Connect(context.Background())
	mock.AssertStateTransition(t, mc, axon.StateDisconnected, axon.StateConnecting)
	mock.AssertStateTransition(t, mc, axon.StateConnecting, axon.StateConnected)

Testing error handling:

	mc := mock.NewMockConn[string](nil)
	defer mc.Close(1000, "")

	mc.InjectReadError(axon.ErrConnectionClosed)
	_, err := mc.Read(context.Background())
	if err != axon.ErrConnectionClosed {
	    t.Fatalf("unexpected error: %v", err)
	}


## Types

### CallbackRecorder
CallbackRecorder records callback invocations for testing

#### Example Usage

```go
// Create a new CallbackRecorder
callbackrecorder := CallbackRecorder{

}
```

#### Type Definition

```go
type CallbackRecorder struct {
}
```

### Constructor Functions

### NewCallbackRecorder

NewCallbackRecorder creates a new callback recorder

```go
func NewCallbackRecorder() *CallbackRecorder
```

**Parameters:**
  None

**Returns:**
- *CallbackRecorder

## Methods

### AssertCallCount

AssertCallCount asserts the exact number of callback invocations

```go
func (*CallbackRecorder) AssertCallCount(t *testing.T, expected int)
```

**Parameters:**
- `t` (*testing.T)
- `expected` (int)

**Returns:**
  None

### AssertCalled

AssertCalled asserts a callback was invoked

```go
func (*CallbackRecorder) AssertCalled(t *testing.T)
```

**Parameters:**
- `t` (*testing.T)

**Returns:**
  None

### AssertNotCalled

AssertNotCalled asserts a callback was not invoked

```go
func (*CallbackRecorder) AssertNotCalled(t *testing.T)
```

**Parameters:**
- `t` (*testing.T)

**Returns:**
  None

### Called

Called returns the number of times the callback was called

```go
func (*CallbackRecorder) Called() int
```

**Parameters:**
  None

**Returns:**
- int

### Record

Record records a callback invocation

```go
func (*CallbackRecorder) Record()
```

**Parameters:**
  None

**Returns:**
  None

### CommonErrors
ErrorInjectionHelpers provides utilities for error injection in tests. Error injection methods are available directly on MockConn[T] and MockClient[T]: mc := mock.NewMockConn[string](nil) mc.InjectReadError(axon.ErrConnectionClosed) _, err := mc.Read(ctx) // err will be axon.ErrConnectionClosed CommonErrors provides frequently-used errors for injection

#### Example Usage

```go
// Create a new CommonErrors
commonerrors := CommonErrors{
    ConnectionClosed: error{},
    ContextCanceled: error{},
    QueueFull: error{},
    InvalidState: error{},
}
```

#### Type Definition

```go
type CommonErrors struct {
    ConnectionClosed error
    ContextCanceled error
    QueueFull error
    InvalidState error
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ConnectionClosed | `error` |  |
| ContextCanceled | `error` |  |
| QueueFull | `error` |  |
| InvalidState | `error` |  |

### Constructor Functions

### DefaultErrors

DefaultErrors returns a set of common errors for injection

```go
func DefaultErrors() *CommonErrors
```

**Parameters:**
  None

**Returns:**
- *CommonErrors

### ErrorMatcher
ErrorMatcher is a function that matches errors for assertions

#### Example Usage

```go
// Example usage of ErrorMatcher
var value ErrorMatcher
// Initialize with appropriate value
```

#### Type Definition

```go
type ErrorMatcher func(error) bool
```

### ErrorSequence
ErrorSequence allows injecting different errors for successive operations

#### Example Usage

```go
// Create a new ErrorSequence
errorsequence := ErrorSequence{

}
```

#### Type Definition

```go
type ErrorSequence struct {
}
```

### Constructor Functions

### NewErrorSequence

NewErrorSequence creates a new error sequence

```go
func NewErrorSequence(errs ...error) *ErrorSequence
```

**Parameters:**
- `errs` (...error)

**Returns:**
- *ErrorSequence

## Methods

### Next

Next returns the next error in the sequence, or nil if exhausted

```go
func (*ErrorSequence) Next() error
```

**Parameters:**
  None

**Returns:**
- error

### Reset

Reset resets the sequence to the beginning

```go
func (*ErrorSequence) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### MessageMatcher
MessageMatcher is a function that matches messages for assertions

#### Example Usage

```go
// Example usage of MessageMatcher
var value MessageMatcher
// Initialize with appropriate value
```

#### Type Definition

```go
type MessageMatcher func(T) bool
```

### MockClient
MockClient provides an in-memory mock of axon.Client[T] for testing

#### Example Usage

```go
// Create a new MockClient
mockclient := MockClient{

}
```

#### Type Definition

```go
type MockClient struct {
}
```

### Constructor Functions

### NewMockClient

NewMockClient creates a new mock client

```go
func NewMockClient(opts *MockClientOptions) **ast.IndexExpr
```

**Parameters:**
- `opts` (*MockClientOptions)

**Returns:**
- **ast.IndexExpr

## Methods

### ClearInjectedErrors

ClearInjectedErrors clears all injected errors

```go
func (**ast.IndexExpr) ClearInjectedErrors()
```

**Parameters:**
  None

**Returns:**
  None

### ClearRecorded

ClearRecorded clears all recorded data

```go
func (**ast.IndexExpr) ClearRecorded()
```

**Parameters:**
  None

**Returns:**
  None

### Close

Close closes the mock client

```go
func (**ast.IndexExpr) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Connect

Connect simulates establishing a connection

```go
func (**ast.IndexExpr) Connect(ctx context.Context) error
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- error

### ConnectWithReadLoop

ConnectWithReadLoop connects and starts message delivery

```go
func (**ast.IndexExpr) ConnectWithReadLoop(ctx context.Context) error
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- error

### DrainWritten

DrainWritten retrieves all written messages

```go
func (**ast.IndexExpr) DrainWritten() []T
```

**Parameters:**
  None

**Returns:**
- []T

### ForceState

ForceState forces a specific state (use with caution in tests)

```go
func (**ast.IndexExpr) ForceState(state axon.ConnectionState)
```

**Parameters:**
- `state` (axon.ConnectionState)

**Returns:**
  None

### InjectConnectError

InjectConnectError injects an error for the next Connect operation

```go
func (**ast.IndexExpr) InjectConnectError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### InjectMessage

InjectMessage simulates receiving a message

```go
func (**ast.IndexExpr) InjectMessage(msg T) error
```

**Parameters:**
- `msg` (T)

**Returns:**
- error

### InjectReadError

InjectReadError injects an error for the next Read operation

```go
func (**ast.IndexExpr) InjectReadError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### InjectWriteError

InjectWriteError injects an error for the next Write operation

```go
func (**ast.IndexExpr) InjectWriteError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

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

OnConnect sets the connect callback

```go
func (**ast.IndexExpr) OnConnect(fn func())
```

**Parameters:**
- `fn` (func())

**Returns:**
  None

### OnDisconnect

OnDisconnect sets the disconnect callback

```go
func (**ast.IndexExpr) OnDisconnect(fn func(error))
```

**Parameters:**
- `fn` (func(error))

**Returns:**
  None

### OnMessage

OnMessage sets the message callback

```go
func (**ast.IndexExpr) OnMessage(fn func(T))
```

**Parameters:**
- `fn` (func(T))

**Returns:**
  None

### OnStateChange

OnStateChange sets the state change callback

```go
func (**ast.IndexExpr) OnStateChange(fn func(axon.StateChange))
```

**Parameters:**
- `fn` (func(axon.StateChange))

**Returns:**
  None

### QueueStats

QueueStats returns queue statistics

```go
func (**ast.IndexExpr) QueueStats() axon.MessageQueueStats
```

**Parameters:**
  None

**Returns:**
- axon.MessageQueueStats

### Read

Read reads a message from the client

```go
func (**ast.IndexExpr) Read(ctx context.Context) (T, error)
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- T
- error

### ReceivedMessages

ReceivedMessages returns recorded received messages

```go
func (**ast.IndexExpr) ReceivedMessages() []T
```

**Parameters:**
  None

**Returns:**
- []T

### SessionID

SessionID returns the session ID

```go
func (**ast.IndexExpr) SessionID() string
```

**Parameters:**
  None

**Returns:**
- string

### SetSessionID

SetSessionID sets the session ID

```go
func (**ast.IndexExpr) SetSessionID(id string)
```

**Parameters:**
- `id` (string)

**Returns:**
  None

### SimulateDisconnect

SimulateDisconnect forces a disconnection with optional error

```go
func (**ast.IndexExpr) SimulateDisconnect(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### SimulateReconnecting

SimulateReconnecting transitions to reconnecting state

```go
func (**ast.IndexExpr) SimulateReconnecting(attempt int)
```

**Parameters:**
- `attempt` (int)

**Returns:**
  None

### State

State returns the current connection state

```go
func (**ast.IndexExpr) State() axon.ConnectionState
```

**Parameters:**
  None

**Returns:**
- axon.ConnectionState

### StateChanges

StateChanges returns recorded state changes

```go
func (**ast.IndexExpr) StateChanges() []axon.StateChange
```

**Parameters:**
  None

**Returns:**
- []axon.StateChange

### Write

Write writes a message to the client

```go
func (**ast.IndexExpr) Write(ctx context.Context, msg T) error
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (T)

**Returns:**
- error

### WrittenMessages

WrittenMessages returns recorded written messages

```go
func (**ast.IndexExpr) WrittenMessages() []T
```

**Parameters:**
  None

**Returns:**
- []T

### MockClientOptions
MockClientOptions configures MockClient behavior

#### Example Usage

```go
// Create a new MockClientOptions
mockclientoptions := MockClientOptions{
    BufferSize: 42,
    InitialState: /* value */,
    ConnectDelay: /* value */,
    ReadDelay: /* value */,
    WriteDelay: /* value */,
    RecordMessages: true,
    RecordStateChanges: true,
    QueueSize: 42,
}
```

#### Type Definition

```go
type MockClientOptions struct {
    BufferSize int
    InitialState axon.ConnectionState
    ConnectDelay time.Duration
    ReadDelay time.Duration
    WriteDelay time.Duration
    RecordMessages bool
    RecordStateChanges bool
    QueueSize int
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| BufferSize | `int` |  |
| InitialState | `axon.ConnectionState` |  |
| ConnectDelay | `time.Duration` |  |
| ReadDelay | `time.Duration` |  |
| WriteDelay | `time.Duration` |  |
| RecordMessages | `bool` |  |
| RecordStateChanges | `bool` |  |
| QueueSize | `int` |  |

### MockConn
MockConn provides an in-memory mock of axon.Conn[T] for testing

#### Example Usage

```go
// Create a new MockConn
mockconn := MockConn{

}
```

#### Type Definition

```go
type MockConn struct {
}
```

### Constructor Functions

### NewMockConn

NewMockConn creates a new mock connection

```go
func NewMockConn(opts *MockConnOptions) **ast.IndexExpr
```

**Parameters:**
- `opts` (*MockConnOptions)

**Returns:**
- **ast.IndexExpr

## Methods

### ClearInjectedErrors

ClearInjectedErrors clears all injected errors

```go
func (**ast.IndexExpr) ClearInjectedErrors()
```

**Parameters:**
  None

**Returns:**
  None

### ClearRecorded

ClearRecorded clears all recorded messages

```go
func (**ast.IndexExpr) ClearRecorded()
```

**Parameters:**
  None

**Returns:**
  None

### Close

Close closes the mock connection

```go
func (**ast.IndexExpr) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### CloseCode

CloseCode returns the close code

```go
func (**ast.IndexExpr) CloseCode() int
```

**Parameters:**
  None

**Returns:**
- int

### CloseReason

CloseReason returns the close reason

```go
func (**ast.IndexExpr) CloseReason() string
```

**Parameters:**
  None

**Returns:**
- string

### DrainWritten

DrainWritten retrieves all written messages

```go
func (**ast.IndexExpr) DrainWritten() []T
```

**Parameters:**
  None

**Returns:**
- []T

### InjectCloseError

InjectCloseError injects an error for the next Close operation

```go
func (**ast.IndexExpr) InjectCloseError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### InjectMessage

InjectMessage simulates receiving a message (test helper)

```go
func (**ast.IndexExpr) InjectMessage(msg T) error
```

**Parameters:**
- `msg` (T)

**Returns:**
- error

### InjectReadError

InjectReadError injects an error for the next Read operation

```go
func (**ast.IndexExpr) InjectReadError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### InjectWriteError

InjectWriteError injects an error for the next Write operation

```go
func (**ast.IndexExpr) InjectWriteError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### IsClosed

IsClosed returns true if the connection is closed

```go
func (**ast.IndexExpr) IsClosed() bool
```

**Parameters:**
  None

**Returns:**
- bool

### Read

Read reads a message from the mock connection

```go
func (**ast.IndexExpr) Read(ctx context.Context) (T, error)
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- T
- error

### ReadMessages

ReadMessages returns recorded read messages

```go
func (**ast.IndexExpr) ReadMessages() []T
```

**Parameters:**
  None

**Returns:**
- []T

### Write

Write writes a message to the mock connection

```go
func (**ast.IndexExpr) Write(ctx context.Context, msg T) error
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (T)

**Returns:**
- error

### WrittenMessages

WrittenMessages returns recorded written messages

```go
func (**ast.IndexExpr) WrittenMessages() []T
```

**Parameters:**
  None

**Returns:**
- []T

### MockConnOptions
MockConnOptions configures MockConn behavior

#### Example Usage

```go
// Create a new MockConnOptions
mockconnoptions := MockConnOptions{
    BufferSize: 42,
    ReadDelay: /* value */,
    WriteDelay: /* value */,
    RecordMessages: true,
}
```

#### Type Definition

```go
type MockConnOptions struct {
    BufferSize int
    ReadDelay time.Duration
    WriteDelay time.Duration
    RecordMessages bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| BufferSize | `int` |  |
| ReadDelay | `time.Duration` |  |
| WriteDelay | `time.Duration` |  |
| RecordMessages | `bool` |  |

## Functions

### AssertClientWritten
AssertClientWritten asserts that a message was written by the client

```go
func AssertClientWritten(t *testing.T, mc **ast.IndexExpr, expected T)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `T` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertClientWritten
result := AssertClientWritten(/* parameters */)
```

### AssertClientWrittenCount
AssertClientWrittenCount asserts the number of written messages by the client

```go
func AssertClientWrittenCount(t *testing.T, mc **ast.IndexExpr, expected int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertClientWrittenCount
result := AssertClientWrittenCount(/* parameters */)
```

### AssertClientWrittenMatches
AssertClientWrittenMatches asserts a message matching the predicate was written by the client

```go
func AssertClientWrittenMatches(t *testing.T, mc **ast.IndexExpr, matcher *ast.IndexExpr)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `matcher` | `*ast.IndexExpr` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertClientWrittenMatches
result := AssertClientWrittenMatches(/* parameters */)
```

### AssertClosed
AssertClosed asserts the client is closed

```go
func AssertClosed(t *testing.T, mc **ast.IndexExpr)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertClosed
result := AssertClosed(/* parameters */)
```

### AssertConnected
AssertConnected asserts the client is connected

```go
func AssertConnected(t *testing.T, mc **ast.IndexExpr)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertConnected
result := AssertConnected(/* parameters */)
```

### AssertDisconnected
AssertDisconnected asserts the client is disconnected

```go
func AssertDisconnected(t *testing.T, mc **ast.IndexExpr)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertDisconnected
result := AssertDisconnected(/* parameters */)
```

### AssertLastStateChange
AssertLastStateChange asserts the last state change matches expectations

```go
func AssertLastStateChange(t *testing.T, mc **ast.IndexExpr, from, to axon.ConnectionState)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `from` | `axon.ConnectionState` | |
| `to` | `axon.ConnectionState` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertLastStateChange
result := AssertLastStateChange(/* parameters */)
```

### AssertReceivedCount
AssertReceivedCount asserts the number of received messages

```go
func AssertReceivedCount(t *testing.T, mc **ast.IndexExpr, expected int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertReceivedCount
result := AssertReceivedCount(/* parameters */)
```

### AssertState
AssertState asserts the client is in the expected state

```go
func AssertState(t *testing.T, mc **ast.IndexExpr, expected axon.ConnectionState)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `axon.ConnectionState` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertState
result := AssertState(/* parameters */)
```

### AssertStateChangeCallback
AssertStateChangeCallback asserts state changes via callback

```go
func AssertStateChangeCallback(t *testing.T, mc **ast.IndexExpr, from, to axon.ConnectionState)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `from` | `axon.ConnectionState` | |
| `to` | `axon.ConnectionState` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertStateChangeCallback
result := AssertStateChangeCallback(/* parameters */)
```

### AssertStateChangeCount
AssertStateChangeCount asserts the number of state changes

```go
func AssertStateChangeCount(t *testing.T, mc **ast.IndexExpr, expected int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertStateChangeCount
result := AssertStateChangeCount(/* parameters */)
```

### AssertStateTransition
AssertStateTransition asserts a state transition occurred

```go
func AssertStateTransition(t *testing.T, mc **ast.IndexExpr, from, to axon.ConnectionState)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `from` | `axon.ConnectionState` | |
| `to` | `axon.ConnectionState` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertStateTransition
result := AssertStateTransition(/* parameters */)
```

### AssertWritten
AssertWritten asserts that a message was written

```go
func AssertWritten(t *testing.T, mc **ast.IndexExpr, expected T)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `T` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertWritten
result := AssertWritten(/* parameters */)
```

### AssertWrittenCount
AssertWrittenCount asserts the number of written messages

```go
func AssertWrittenCount(t *testing.T, mc **ast.IndexExpr, expected int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertWrittenCount
result := AssertWrittenCount(/* parameters */)
```

### AssertWrittenMatches
AssertWrittenMatches asserts a message matching the predicate was written

```go
func AssertWrittenMatches(t *testing.T, mc **ast.IndexExpr, matcher *ast.IndexExpr)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `matcher` | `*ast.IndexExpr` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertWrittenMatches
result := AssertWrittenMatches(/* parameters */)
```

### IsConnectionClosedError
IsConnectionClosedError returns true if the error is a connection closed error

```go
func IsConnectionClosedError(err error) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `err` | `error` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of IsConnectionClosedError
result := IsConnectionClosedError(/* parameters */)
```

### IsContextCancelledError
IsContextCancelledError returns true if the error is a context canceled error

```go
func IsContextCancelledError(err error) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `err` | `error` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of IsContextCancelledError
result := IsContextCancelledError(/* parameters */)
```

### IsQueueFullError
IsQueueFullError returns true if the error is a queue full error

```go
func IsQueueFullError(err error) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `err` | `error` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of IsQueueFullError
result := IsQueueFullError(/* parameters */)
```

### NoOpCallback
NoOpCallback returns a no-op callback for testing without side effects

```go
func NoOpCallback() func(T)
```

**Parameters:**
None

**Returns:**
| Type | Description |
|------|-------------|
| `func(T)` | |

**Example:**

```go
// Example usage of NoOpCallback
result := NoOpCallback(/* parameters */)
```

### PanicOnError
PanicOnError returns a function that panics on error (for use in callbacks)

```go
func PanicOnError() func(error)
```

**Parameters:**
None

**Returns:**
| Type | Description |
|------|-------------|
| `func(error)` | |

**Example:**

```go
// Example usage of PanicOnError
result := PanicOnError(/* parameters */)
```

### WaitForClientMessage
WaitForClientMessage waits for a message to be written by the client

```go
func WaitForClientMessage(t *testing.T, mc **ast.IndexExpr, expected T, timeout time.Duration)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `T` | |
| `timeout` | `time.Duration` | |

**Returns:**
None

**Example:**

```go
// Example usage of WaitForClientMessage
result := WaitForClientMessage(/* parameters */)
```

### WaitForMessage
WaitForMessage waits for a message to be written

```go
func WaitForMessage(t *testing.T, mc **ast.IndexExpr, expected T, timeout time.Duration)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `T` | |
| `timeout` | `time.Duration` | |

**Returns:**
None

**Example:**

```go
// Example usage of WaitForMessage
result := WaitForMessage(/* parameters */)
```

### WaitForState
WaitForState waits for the client to reach the expected state

```go
func WaitForState(t *testing.T, mc **ast.IndexExpr, expected axon.ConnectionState, timeout time.Duration)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `mc` | `**ast.IndexExpr` | |
| `expected` | `axon.ConnectionState` | |
| `timeout` | `time.Duration` | |

**Returns:**
None

**Example:**

```go
// Example usage of WaitForState
result := WaitForState(/* parameters */)
```

## External Links

- [Package Overview](../packages/mock.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/axon/mock)
- [Source Code](https://github.com/kolosys/axon/tree/main/mock)
