package mock_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kolosys/axon"
	"github.com/kolosys/axon/mock"
)

func TestMockClient_Connect(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	if mc.State() != axon.StateDisconnected {
		t.Errorf("State() = %v, want StateDisconnected", mc.State())
	}

	err := mc.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	if mc.State() != axon.StateConnected {
		t.Errorf("State() = %v, want StateConnected", mc.State())
	}

	if !mc.IsConnected() {
		t.Error("expected IsConnected() to be true")
	}
}

func TestMockClient_Callbacks(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	connectRec := mock.NewCallbackRecorder()
	disconnectRec := mock.NewCallbackRecorder()
	stateChangeRec := mock.NewCallbackRecorder()

	mc.OnConnect(connectRec.Record)
	mc.OnDisconnect(func(error) {
		disconnectRec.Record()
	})
	mc.OnStateChange(func(axon.StateChange) {
		stateChangeRec.Record()
	})

	err := mc.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	connectRec.AssertCalled(t)

	mc.SimulateDisconnect(axon.ErrConnectionClosed)
	disconnectRec.AssertCalled(t)
	stateChangeRec.AssertCalled(t)
}

func TestMockClient_StateTransitions(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.Connect(context.Background())

	changes := mc.StateChanges()
	if len(changes) < 2 {
		t.Fatalf("expected at least 2 state changes, got %d", len(changes))
	}

	if changes[0].From != axon.StateDisconnected || changes[0].To != axon.StateConnecting {
		t.Errorf("first transition = %v->%v, want StateDisconnected->StateConnecting", changes[0].From, changes[0].To)
	}

	if changes[1].From != axon.StateConnecting || changes[1].To != axon.StateConnected {
		t.Errorf("second transition = %v->%v, want StateConnecting->StateConnected", changes[1].From, changes[1].To)
	}
}

func TestMockClient_MessageRecording(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.Connect(context.Background())

	ctx := context.Background()

	mc.Write(ctx, "msg1")
	mc.Write(ctx, "msg2")
	mc.Write(ctx, "msg3")

	written := mc.WrittenMessages()
	if len(written) != 3 {
		t.Fatalf("expected 3 written messages, got %d", len(written))
	}

	expected := []string{"msg1", "msg2", "msg3"}
	for i, msg := range written {
		if msg != expected[i] {
			t.Errorf("WrittenMessages()[%d] = %q, want %q", i, msg, expected[i])
		}
	}
}

func TestMockClient_SessionID(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	if mc.SessionID() != "" {
		t.Errorf("SessionID() = %q, want empty", mc.SessionID())
	}

	mc.SetSessionID("test-123")

	if mc.SessionID() != "test-123" {
		t.Errorf("SessionID() = %q, want test-123", mc.SessionID())
	}
}

func TestMockClient_ConnectError(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	testErr := errors.New("connection failed")
	mc.InjectConnectError(testErr)

	err := mc.Connect(context.Background())
	if err != testErr {
		t.Errorf("Connect() error = %v, want %v", err, testErr)
	}
}

func TestMockClient_ReadWriteErrors(t *testing.T) {
	tests := []struct {
		name      string
		injectFn  func(*mock.MockClient[string])
		operation string
	}{
		{
			name: "read error",
			injectFn: func(mc *mock.MockClient[string]) {
				mc.InjectReadError(axon.ErrConnectionClosed)
			},
			operation: "read",
		},
		{
			name: "write error",
			injectFn: func(mc *mock.MockClient[string]) {
				mc.InjectWriteError(axon.ErrConnectionClosed)
			},
			operation: "write",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := mock.NewMockClient[string](nil)
			defer mc.Close()

			mc.Connect(context.Background())
			tt.injectFn(mc)

			ctx := context.Background()
			var err error

			switch tt.operation {
			case "read":
				_, err = mc.Read(ctx)
			case "write":
				err = mc.Write(ctx, "test")
			}

			if err != axon.ErrConnectionClosed {
				t.Errorf("expected ErrConnectionClosed, got %v", err)
			}
		})
	}
}

func TestMockClient_SimulateDisconnect(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.Connect(context.Background())

	if mc.State() != axon.StateConnected {
		t.Fatalf("expected to be connected")
	}

	mc.SimulateDisconnect(axon.ErrConnectionClosed)

	if mc.State() != axon.StateDisconnected {
		t.Errorf("State() = %v, want StateDisconnected", mc.State())
	}
}

func TestMockClient_SimulateReconnecting(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.Connect(context.Background())
	mc.SimulateDisconnect(nil)

	mc.SimulateReconnecting(1)

	if mc.State() != axon.StateReconnecting {
		t.Errorf("State() = %v, want StateReconnecting", mc.State())
	}
}

func TestMockClient_ForceState(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.ForceState(axon.StateConnected)

	if mc.State() != axon.StateConnected {
		t.Errorf("State() = %v, want StateConnected", mc.State())
	}

	if !mc.IsConnected() {
		t.Error("expected IsConnected() to be true")
	}
}

func TestMockClient_ClearRecorded(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.Connect(context.Background())

	ctx := context.Background()
	mc.Write(ctx, "msg1")
	mc.Write(ctx, "msg2")

	written := mc.WrittenMessages()
	if len(written) != 2 {
		t.Errorf("expected 2 messages before clear, got %d", len(written))
	}

	mc.ClearRecorded()

	written = mc.WrittenMessages()
	if len(written) != 0 {
		t.Errorf("expected 0 messages after clear, got %d", len(written))
	}

	changes := mc.StateChanges()
	if len(changes) != 0 {
		t.Errorf("expected 0 state changes after clear, got %d", len(changes))
	}
}

func TestMockClient_NoRecording(t *testing.T) {
	opts := &mock.MockClientOptions{
		RecordMessages:     false,
		RecordStateChanges: false,
	}
	mc := mock.NewMockClient[string](opts)
	defer mc.Close()

	mc.Connect(context.Background())

	ctx := context.Background()
	mc.Write(ctx, "msg1")

	written := mc.WrittenMessages()
	if written != nil {
		t.Errorf("WrittenMessages() = %v, want nil", written)
	}

	changes := mc.StateChanges()
	if changes != nil {
		t.Errorf("StateChanges() = %v, want nil", changes)
	}
}

func TestMockClient_ContextCancellation(t *testing.T) {
	opts := &mock.MockClientOptions{
		ConnectDelay: 100 * time.Millisecond,
	}
	mc := mock.NewMockClient[string](opts)
	defer mc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := mc.Connect(ctx)
	if err != axon.ErrContextCanceled {
		t.Errorf("Connect() error = %v, want ErrContextCanceled", err)
	}
}

func TestMockClient_WriteWhenDisconnected(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	err := mc.Write(context.Background(), "test")
	if err != axon.ErrConnectionClosed {
		t.Errorf("Write() error = %v, want ErrConnectionClosed", err)
	}
}

func TestMockClient_ReadWhenDisconnected(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	_, err := mc.Read(context.Background())
	if err != axon.ErrConnectionClosed {
		t.Errorf("Read() error = %v, want ErrConnectionClosed", err)
	}
}

func TestMockClient_ConnectDelay(t *testing.T) {
	opts := &mock.MockClientOptions{
		ConnectDelay: 50 * time.Millisecond,
	}
	mc := mock.NewMockClient[string](opts)
	defer mc.Close()

	start := time.Now()
	mc.Connect(context.Background())
	elapsed := time.Since(start)

	if elapsed < 50*time.Millisecond {
		t.Errorf("Connect() delay was %v, want >= 50ms", elapsed)
	}
}

func TestMockClient_QueueStats(t *testing.T) {
	opts := &mock.MockClientOptions{
		QueueSize: 256,
	}
	mc := mock.NewMockClient[string](opts)
	defer mc.Close()

	stats := mc.QueueStats()
	if stats.MaxSize != 256 {
		t.Errorf("QueueStats().MaxSize = %d, want 256", stats.MaxSize)
	}
}

func TestMockClient_InjectMessage(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.Connect(context.Background())

	err := mc.InjectMessage("hello")
	if err != nil {
		t.Fatalf("InjectMessage() error = %v", err)
	}

	msg, err := mc.Read(context.Background())
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if msg != "hello" {
		t.Errorf("Read() = %q, want hello", msg)
	}
}

func TestMockClient_DrainWritten(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.Connect(context.Background())

	ctx := context.Background()
	mc.Write(ctx, "msg1")
	mc.Write(ctx, "msg2")

	msgs := mc.DrainWritten()
	if len(msgs) != 2 {
		t.Errorf("expected 2 drained messages, got %d", len(msgs))
	}
}

func TestMockClient_StateChangeCallback(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	stateChangeCount := 0
	mc.OnStateChange(func(change axon.StateChange) {
		stateChangeCount++
	})

	mc.Connect(context.Background())

	if stateChangeCount < 2 {
		t.Errorf("expected at least 2 state change callbacks, got %d", stateChangeCount)
	}
}

func TestMockClient_ConnectAlreadyConnected(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	err := mc.Connect(context.Background())
	if err != nil {
		t.Fatalf("first Connect() error = %v", err)
	}

	err = mc.Connect(context.Background())
	if err != nil {
		t.Fatalf("second Connect() error = %v (should be nil when already connected)", err)
	}
}

func TestMockClient_ClearInjectedErrors(t *testing.T) {
	mc := mock.NewMockClient[string](nil)
	defer mc.Close()

	mc.ForceState(axon.StateConnected)
	mc.InjectReadError(axon.ErrConnectionClosed)
	mc.InjectWriteError(axon.ErrConnectionClosed)
	mc.ClearInjectedErrors()

	ctx := context.Background()

	mc.InjectMessage("test")
	_, err := mc.Read(ctx)
	if err != nil {
		t.Errorf("Read() error = %v, want nil after clearing", err)
	}

	err = mc.Write(ctx, "test")
	if err != nil {
		t.Errorf("Write() error = %v, want nil after clearing", err)
	}
}
