package mock_test

import (
	"context"
	"testing"
	"time"

	"github.com/kolosys/axon"
	"github.com/kolosys/axon/mock"
)

func TestMockConn_ReadWrite(t *testing.T) {
	mc := mock.NewMockConn[string](nil)
	defer mc.Close(1000, "test done")

	ctx := context.Background()

	err := mc.InjectMessage("hello")
	if err != nil {
		t.Fatalf("InjectMessage() error = %v", err)
	}

	msg, err := mc.Read(ctx)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if msg != "hello" {
		t.Errorf("Read() = %q, want %q", msg, "hello")
	}

	err = mc.Write(ctx, "world")
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	written := mc.WrittenMessages()
	if len(written) != 1 || written[0] != "world" {
		t.Errorf("WrittenMessages() = %v, want [world]", written)
	}
}

func TestMockConn_ErrorInjection(t *testing.T) {
	tests := []struct {
		name      string
		injectFn  func(*mock.MockConn[string])
		operation string
	}{
		{
			name: "read error",
			injectFn: func(mc *mock.MockConn[string]) {
				mc.InjectReadError(axon.ErrConnectionClosed)
			},
			operation: "read",
		},
		{
			name: "write error",
			injectFn: func(mc *mock.MockConn[string]) {
				mc.InjectWriteError(axon.ErrConnectionClosed)
			},
			operation: "write",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := mock.NewMockConn[string](nil)
			defer mc.Close(1000, "")

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

func TestMockConn_Close(t *testing.T) {
	mc := mock.NewMockConn[string](nil)

	if mc.IsClosed() {
		t.Error("expected IsClosed() to be false")
	}

	err := mc.Close(1000, "normal")
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if !mc.IsClosed() {
		t.Error("expected IsClosed() to be true")
	}

	if code := mc.CloseCode(); code != 1000 {
		t.Errorf("CloseCode() = %d, want 1000", code)
	}

	if reason := mc.CloseReason(); reason != "normal" {
		t.Errorf("CloseReason() = %q, want %q", reason, "normal")
	}
}

func TestMockConn_ContextCancellation(t *testing.T) {
	opts := &mock.MockConnOptions{
		ReadDelay:  100 * time.Millisecond,
		WriteDelay: 100 * time.Millisecond,
	}
	mc := mock.NewMockConn[string](opts)
	defer mc.Close(1000, "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := mc.Read(ctx)
	if err != axon.ErrContextCanceled {
		t.Errorf("Read() error = %v, want ErrContextCanceled", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel2()
	err = mc.Write(ctx2, "test")
	if err != axon.ErrContextCanceled {
		t.Errorf("Write() error = %v, want ErrContextCanceled", err)
	}
}

func TestMockConn_ClosedConnection(t *testing.T) {
	mc := mock.NewMockConn[string](nil)
	mc.Close(1000, "")

	ctx := context.Background()

	_, err := mc.Read(ctx)
	if err != axon.ErrConnectionClosed {
		t.Errorf("Read() error = %v, want ErrConnectionClosed", err)
	}

	err = mc.Write(ctx, "test")
	if err != axon.ErrConnectionClosed {
		t.Errorf("Write() error = %v, want ErrConnectionClosed", err)
	}
}

func TestMockConn_MessageRecording(t *testing.T) {
	mc := mock.NewMockConn[string](nil)
	defer mc.Close(1000, "")

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

func TestMockConn_MultipleTypes(t *testing.T) {
	type Message struct {
		ID   int
		Text string
	}

	mc := mock.NewMockConn[Message](nil)
	defer mc.Close(1000, "")

	ctx := context.Background()

	msg := Message{ID: 1, Text: "hello"}
	err := mc.Write(ctx, msg)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	written := mc.WrittenMessages()
	if len(written) != 1 || written[0].ID != 1 {
		t.Errorf("WrittenMessages() = %v, want [%v]", written, msg)
	}
}

func TestMockConn_ClearRecorded(t *testing.T) {
	mc := mock.NewMockConn[string](nil)
	defer mc.Close(1000, "")

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
}

func TestMockConn_DrainWritten(t *testing.T) {
	mc := mock.NewMockConn[string](nil)
	defer mc.Close(1000, "")

	ctx := context.Background()
	mc.Write(ctx, "msg1")
	mc.Write(ctx, "msg2")

	msgs := mc.DrainWritten()
	if len(msgs) != 2 {
		t.Errorf("expected 2 drained messages, got %d", len(msgs))
	}
}

func TestMockConn_InjectQueueFull(t *testing.T) {
	opts := &mock.MockConnOptions{
		BufferSize: 1,
	}
	mc := mock.NewMockConn[string](opts)
	defer mc.Close(1000, "")

	err := mc.InjectMessage("msg1")
	if err != nil {
		t.Fatalf("first InjectMessage() error = %v", err)
	}

	err = mc.InjectMessage("msg2")
	if err != axon.ErrQueueFull {
		t.Errorf("second InjectMessage() error = %v, want ErrQueueFull", err)
	}
}

func TestMockConn_ConcurrentReadWrite(t *testing.T) {
	mc := mock.NewMockConn[string](nil)
	defer mc.Close(1000, "")

	ctx := context.Background()
	done := make(chan bool)

	go func() {
		mc.Write(ctx, "msg1")
		mc.Write(ctx, "msg2")
		done <- true
	}()

	mc.InjectMessage("received1")
	mc.InjectMessage("received2")

	msg, _ := mc.Read(ctx)
	if msg != "received1" {
		t.Errorf("Read() = %q, want received1", msg)
	}

	<-done
}

func TestMockConn_NoRecording(t *testing.T) {
	opts := &mock.MockConnOptions{
		RecordMessages: false,
	}
	mc := mock.NewMockConn[string](opts)
	defer mc.Close(1000, "")

	ctx := context.Background()
	mc.Write(ctx, "msg1")

	written := mc.WrittenMessages()
	if written != nil {
		t.Errorf("WrittenMessages() = %v, want nil", written)
	}
}

func TestMockConn_CloseError(t *testing.T) {
	mc := mock.NewMockConn[string](nil)

	testErr := axon.NewCloseError(4000, "test error")
	mc.InjectCloseError(testErr)

	err := mc.Close(1000, "normal")
	if err != testErr {
		t.Errorf("Close() error = %v, want %v", err, testErr)
	}
}

func TestMockConn_ReadDelay(t *testing.T) {
	opts := &mock.MockConnOptions{
		ReadDelay: 50 * time.Millisecond,
	}
	mc := mock.NewMockConn[string](opts)
	defer mc.Close(1000, "")

	mc.InjectMessage("hello")

	start := time.Now()
	mc.Read(context.Background())
	elapsed := time.Since(start)

	if elapsed < 50*time.Millisecond {
		t.Errorf("Read() delay was %v, want >= 50ms", elapsed)
	}
}

func TestMockConn_WriteDelay(t *testing.T) {
	opts := &mock.MockConnOptions{
		WriteDelay: 50 * time.Millisecond,
	}
	mc := mock.NewMockConn[string](opts)
	defer mc.Close(1000, "")

	start := time.Now()
	mc.Write(context.Background(), "hello")
	elapsed := time.Since(start)

	if elapsed < 50*time.Millisecond {
		t.Errorf("Write() delay was %v, want >= 50ms", elapsed)
	}
}

func TestMockConn_ClearInjectedErrors(t *testing.T) {
	mc := mock.NewMockConn[string](nil)
	defer mc.Close(1000, "")

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
