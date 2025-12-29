// Package mock provides in-memory mocks for testing code that uses the axon WebSocket library.
//
// # Overview
//
// The mock package offers two primary mock types:
//   - MockConn[T]: Mocks the low-level Conn[T] connection for unit testing
//   - MockClient[T]: Mocks the high-level Client[T] with full state machine support
//
// Both types use in-memory channels for message passing and provide comprehensive
// testing utilities including error injection, message recording, and state assertions.
//
// # Thread Safety
//
// All mock types are safe for concurrent use. Message recording, state management,
// and error injection all use appropriate synchronization primitives (atomics,
// RWMutex) for thread-safe operation.
//
// # Usage Example - MockConn[T]
//
// MockConn[T] is useful for testing code that directly handles WebSocket connections:
//
//	func TestWebSocketHandler(t *testing.T) {
//	    mc := mock.NewMockConn[string](nil)
//	    defer mc.Close(1000, "test done")
//
//	    // Inject a message to be read
//	    mc.InjectMessage("hello")
//
//	    // Code under test reads from connection
//	    msg, err := mc.Read(context.Background())
//	    if err != nil {
//	        t.Fatal(err)
//	    }
//
//	    // Assert the message
//	    if msg != "hello" {
//	        t.Errorf("got %q, want %q", msg, "hello")
//	    }
//
//	    // Code under test writes response
//	    mc.Write(context.Background(), "world")
//
//	    // Assert what was written
//	    mock.AssertWritten(t, mc, "world")
//	}
//
// # Usage Example - MockClient[T]
//
// MockClient[T] is useful for testing high-level client behavior including
// state transitions and reconnection logic:
//
//	func TestClientReconnection(t *testing.T) {
//	    mc := mock.NewMockClient[string](nil)
//	    defer mc.Close()
//
//	    // Track callbacks
//	    connectCalled := mock.NewCallbackRecorder()
//	    mc.OnConnect(connectCalled.Record)
//
//	    // Connect the client
//	    err := mc.Connect(context.Background())
//	    if err != nil {
//	        t.Fatal(err)
//	    }
//
//	    // Assert state and callback
//	    mock.AssertConnected(t, mc)
//	    connectCalled.AssertCalled(t)
//
//	    // Simulate disconnection
//	    mc.SimulateDisconnect(axon.ErrConnectionClosed)
//	    mock.AssertDisconnected(t, mc)
//	}
//
// # Error Injection
//
// Both mock types support error injection for testing error handling paths:
//
//	mc := mock.NewMockConn[string](nil)
//	mc.InjectReadError(axon.ErrConnectionClosed)
//
//	_, err := mc.Read(context.Background())
//	// err will be axon.ErrConnectionClosed
//
// Clear errors when done with them:
//
//	mc.ClearInjectedErrors()
//
// # Message Recording
//
// Messages are recorded by default and can be retrieved for assertions:
//
//	mc := mock.NewMockClient[string](nil)
//	mc.Connect(context.Background())
//
//	mc.Write(context.Background(), "msg1")
//	mc.Write(context.Background(), "msg2")
//
//	msgs := mc.WrittenMessages()
//	// msgs == []string{"msg1", "msg2"}
//
// Disable recording for performance-sensitive tests:
//
//	opts := &mock.MockClientOptions{
//	    RecordMessages: false,
//	}
//	mc := mock.NewMockClient[string](opts)
//
// # State Management
//
// MockClient[T] provides full state machine support matching the real Client:
//
//	mc := mock.NewMockClient[string](nil)
//
//	// Track state changes
//	mc.OnStateChange(func(change axon.StateChange) {
//	    fmt.Printf("State: %v -> %v\n", change.From, change.To)
//	})
//
//	mc.Connect(context.Background())
//	// Output: State: disconnected -> connecting
//	//         State: connecting -> connected
//
// # Assertion Helpers
//
// The package provides helper functions for common assertions:
//
//	// Assert connection state
//	mock.AssertConnected(t, mc)
//	mock.AssertDisconnected(t, mc)
//	mock.AssertState(t, mc, axon.StateConnected)
//
//	// Assert state transitions
//	mock.AssertStateTransition(t, mc, axon.StateDisconnected, axon.StateConnecting)
//
//	// Assert messages
//	mock.AssertClientWritten(t, mc, expectedMsg)
//	mock.AssertClientWrittenCount(t, mc, 5)
//
//	// Wait for state or message with timeout
//	mock.WaitForState(t, mc, axon.StateConnected, 5*time.Second)
//	mock.WaitForClientMessage(t, mc, "expected", 1*time.Second)
//
// # Callback Testing
//
// Use CallbackRecorder to verify callbacks are invoked:
//
//	mc := mock.NewMockClient[string](nil)
//
//	connectCalled := mock.NewCallbackRecorder()
//	mc.OnConnect(connectCalled.Record)
//
//	mc.Connect(context.Background())
//	connectCalled.AssertCalled(t)
//	connectCalled.AssertCallCount(t, 1)
//
// # Latency Simulation
//
// Simulate network latency for performance testing:
//
//	opts := &mock.MockConnOptions{
//	    ReadDelay:  50 * time.Millisecond,
//	    WriteDelay: 100 * time.Millisecond,
//	}
//	mc := mock.NewMockConn[string](opts)
//
// # Controlled Disconnections
//
// Simulate various disconnection scenarios:
//
//	mc := mock.NewMockClient[string](nil)
//	mc.Connect(context.Background())
//
//	// Graceful disconnection
//	mc.SimulateDisconnect(nil)
//
//	// Disconnection with error
//	mc.SimulateDisconnect(axon.ErrConnectionClosed)
//
//	// Reconnection attempt
//	mc.SimulateReconnecting(1) // attempt number 1
//
// # Type Safety
//
// All mocks are fully generic and maintain type safety:
//
//	type Message struct {
//	    ID   int
//	    Text string
//	}
//
//	mc := mock.NewMockClient[Message](nil)
//	mc.Connect(context.Background())
//
//	msg := Message{ID: 1, Text: "hello"}
//	mc.Write(context.Background(), msg)
//	// Type-safe, no casting needed
//
// # Best Practices
//
//   - Always defer Close() on mocks to ensure cleanup
//   - Use CallbackRecorder to verify callback behavior
//   - Inject errors to test error handling paths
//   - Use assertion helpers for readable test code
//   - Disable recording in performance-critical tests
//   - Clear injected errors between test cases
//
// # Common Patterns
//
// Testing message flow:
//
//	mc := mock.NewMockConn[string](nil)
//	defer mc.Close(1000, "")
//
//	mc.InjectMessage("server message")
//	msg, _ := mc.Read(context.Background())
//	mock.AssertWritten(t, mc, expectedReply)
//
// Testing state transitions:
//
//	mc := mock.NewMockClient[string](nil)
//	defer mc.Close()
//
//	mc.Connect(context.Background())
//	mock.AssertStateTransition(t, mc, axon.StateDisconnected, axon.StateConnecting)
//	mock.AssertStateTransition(t, mc, axon.StateConnecting, axon.StateConnected)
//
// Testing error handling:
//
//	mc := mock.NewMockConn[string](nil)
//	defer mc.Close(1000, "")
//
//	mc.InjectReadError(axon.ErrConnectionClosed)
//	_, err := mc.Read(context.Background())
//	if err != axon.ErrConnectionClosed {
//	    t.Fatalf("unexpected error: %v", err)
//	}
//
package mock
