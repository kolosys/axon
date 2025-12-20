package axon_test

import (
	"testing"
	"time"

	"github.com/kolosys/axon"
)

func TestMetrics(t *testing.T) {
	metrics := &axon.Metrics{}

	metrics.RecordConnection()
	if metrics.ActiveConnections.Load() != 1 {
		t.Errorf("expected 1 active connection, got %d", metrics.ActiveConnections.Load())
	}

	metrics.RecordRead(1024, 100*time.Microsecond)
	if metrics.MessagesRead.Load() != 1 {
		t.Errorf("expected 1 message read, got %d", metrics.MessagesRead.Load())
	}
	if metrics.BytesRead.Load() != 1024 {
		t.Errorf("expected 1024 bytes read, got %d", metrics.BytesRead.Load())
	}

	metrics.RecordWrite(2048, 200*time.Microsecond)
	if metrics.MessagesWritten.Load() != 1 {
		t.Errorf("expected 1 message written, got %d", metrics.MessagesWritten.Load())
	}
	if metrics.BytesWritten.Load() != 2048 {
		t.Errorf("expected 2048 bytes written, got %d", metrics.BytesWritten.Load())
	}

	metrics.RecordDisconnection()
	if metrics.ActiveConnections.Load() != 0 {
		t.Errorf("expected 0 active connections, got %d", metrics.ActiveConnections.Load())
	}
	if metrics.ClosedConnections.Load() != 1 {
		t.Errorf("expected 1 closed connection, got %d", metrics.ClosedConnections.Load())
	}

	metrics.RecordReadError()
	if metrics.ReadErrors.Load() != 1 {
		t.Errorf("expected 1 read error, got %d", metrics.ReadErrors.Load())
	}

	metrics.RecordWriteError()
	if metrics.WriteErrors.Load() != 1 {
		t.Errorf("expected 1 write error, got %d", metrics.WriteErrors.Load())
	}

	metrics.RecordFrameError()
	if metrics.FrameErrors.Load() != 1 {
		t.Errorf("expected 1 frame error, got %d", metrics.FrameErrors.Load())
	}

	metrics.RecordHandshakeError()
	if metrics.HandshakeErrors.Load() != 1 {
		t.Errorf("expected 1 handshake error, got %d", metrics.HandshakeErrors.Load())
	}
}

func TestMetricsSnapshot(t *testing.T) {
	metrics := &axon.Metrics{}

	metrics.RecordConnection()
	metrics.RecordRead(1024, 100*time.Microsecond)
	metrics.RecordWrite(2048, 200*time.Microsecond)

	snapshot := metrics.GetSnapshot()

	if snapshot.ActiveConnections != 1 {
		t.Errorf("expected 1 active connection, got %d", snapshot.ActiveConnections)
	}
	if snapshot.MessagesRead != 1 {
		t.Errorf("expected 1 message read, got %d", snapshot.MessagesRead)
	}
	if snapshot.MessagesWritten != 1 {
		t.Errorf("expected 1 message written, got %d", snapshot.MessagesWritten)
	}
	if snapshot.BytesRead != 1024 {
		t.Errorf("expected 1024 bytes read, got %d", snapshot.BytesRead)
	}
	if snapshot.BytesWritten != 2048 {
		t.Errorf("expected 2048 bytes written, got %d", snapshot.BytesWritten)
	}
}
