package axon

import (
	"sync/atomic"
	"time"
)

// Metrics tracks WebSocket connection metrics
type Metrics struct {
	// Connection metrics
	ActiveConnections atomic.Int64
	TotalConnections  atomic.Int64
	ClosedConnections atomic.Int64

	// Message metrics
	MessagesRead    atomic.Int64
	MessagesWritten atomic.Int64
	BytesRead       atomic.Int64
	BytesWritten    atomic.Int64

	// Error metrics
	ReadErrors      atomic.Int64
	WriteErrors     atomic.Int64
	FrameErrors     atomic.Int64
	HandshakeErrors atomic.Int64

	// Performance metrics
	ReadLatency  atomic.Int64 // nanoseconds
	WriteLatency atomic.Int64 // nanoseconds
}

// MetricsSnapshot represents a snapshot of metrics at a point in time
type MetricsSnapshot struct {
	ActiveConnections int64
	TotalConnections  int64
	ClosedConnections int64
	MessagesRead      int64
	MessagesWritten   int64
	BytesRead         int64
	BytesWritten      int64
	ReadErrors        int64
	WriteErrors       int64
	FrameErrors       int64
	HandshakeErrors   int64
	AvgReadLatency    time.Duration
	AvgWriteLatency   time.Duration
}

// GetSnapshot returns a snapshot of current metrics
func (m *Metrics) GetSnapshot() MetricsSnapshot {
	readCount := m.MessagesRead.Load()
	writeCount := m.MessagesWritten.Load()

	avgReadLatency := time.Duration(0)
	if readCount > 0 {
		avgReadLatency = time.Duration(m.ReadLatency.Load() / readCount)
	}

	avgWriteLatency := time.Duration(0)
	if writeCount > 0 {
		avgWriteLatency = time.Duration(m.WriteLatency.Load() / writeCount)
	}

	return MetricsSnapshot{
		ActiveConnections: m.ActiveConnections.Load(),
		TotalConnections:  m.TotalConnections.Load(),
		ClosedConnections: m.ClosedConnections.Load(),
		MessagesRead:      readCount,
		MessagesWritten:   writeCount,
		BytesRead:         m.BytesRead.Load(),
		BytesWritten:      m.BytesWritten.Load(),
		ReadErrors:        m.ReadErrors.Load(),
		WriteErrors:       m.WriteErrors.Load(),
		FrameErrors:       m.FrameErrors.Load(),
		HandshakeErrors:   m.HandshakeErrors.Load(),
		AvgReadLatency:    avgReadLatency,
		AvgWriteLatency:   avgWriteLatency,
	}
}

// RecordRead records a read operation
func (m *Metrics) RecordRead(bytes int, latency time.Duration) {
	m.MessagesRead.Add(1)
	m.BytesRead.Add(int64(bytes))
	m.ReadLatency.Add(latency.Nanoseconds())
}

// RecordWrite records a write operation
func (m *Metrics) RecordWrite(bytes int, latency time.Duration) {
	m.MessagesWritten.Add(1)
	m.BytesWritten.Add(int64(bytes))
	m.WriteLatency.Add(latency.Nanoseconds())
}

// RecordConnection records a new connection
func (m *Metrics) RecordConnection() {
	m.ActiveConnections.Add(1)
	m.TotalConnections.Add(1)
}

// RecordDisconnection records a connection closure
func (m *Metrics) RecordDisconnection() {
	m.ActiveConnections.Add(-1)
	m.ClosedConnections.Add(1)
}

// RecordReadError records a read error
func (m *Metrics) RecordReadError() {
	m.ReadErrors.Add(1)
}

// RecordWriteError records a write error
func (m *Metrics) RecordWriteError() {
	m.WriteErrors.Add(1)
}

// RecordFrameError records a frame parsing error
func (m *Metrics) RecordFrameError() {
	m.FrameErrors.Add(1)
}

// RecordHandshakeError records a handshake error
func (m *Metrics) RecordHandshakeError() {
	m.HandshakeErrors.Add(1)
}

// DefaultMetrics is the default metrics instance
var DefaultMetrics = &Metrics{}
