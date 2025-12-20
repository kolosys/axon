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

	// Reconnection metrics
	ReconnectAttempts  atomic.Int64
	ReconnectSuccesses atomic.Int64
	ReconnectFailures  atomic.Int64

	// Queue metrics
	QueueEnqueued atomic.Int64
	QueueSent     atomic.Int64
	QueueDropped  atomic.Int64

	// Compression metrics
	CompressedMessages   atomic.Int64
	DecompressedMessages atomic.Int64
	CompressionSaved     atomic.Int64 // bytes saved by compression
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

	// Reconnection metrics
	ReconnectAttempts  int64
	ReconnectSuccesses int64
	ReconnectFailures  int64

	// Queue metrics
	QueueEnqueued int64
	QueueSent     int64
	QueueDropped  int64

	// Compression metrics
	CompressedMessages   int64
	DecompressedMessages int64
	CompressionSaved     int64
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
		ActiveConnections:    m.ActiveConnections.Load(),
		TotalConnections:     m.TotalConnections.Load(),
		ClosedConnections:    m.ClosedConnections.Load(),
		MessagesRead:         readCount,
		MessagesWritten:      writeCount,
		BytesRead:            m.BytesRead.Load(),
		BytesWritten:         m.BytesWritten.Load(),
		ReadErrors:           m.ReadErrors.Load(),
		WriteErrors:          m.WriteErrors.Load(),
		FrameErrors:          m.FrameErrors.Load(),
		HandshakeErrors:      m.HandshakeErrors.Load(),
		AvgReadLatency:       avgReadLatency,
		AvgWriteLatency:      avgWriteLatency,
		ReconnectAttempts:    m.ReconnectAttempts.Load(),
		ReconnectSuccesses:   m.ReconnectSuccesses.Load(),
		ReconnectFailures:    m.ReconnectFailures.Load(),
		QueueEnqueued:        m.QueueEnqueued.Load(),
		QueueSent:            m.QueueSent.Load(),
		QueueDropped:         m.QueueDropped.Load(),
		CompressedMessages:   m.CompressedMessages.Load(),
		DecompressedMessages: m.DecompressedMessages.Load(),
		CompressionSaved:     m.CompressionSaved.Load(),
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

// RecordReconnectAttempt records a reconnection attempt
func (m *Metrics) RecordReconnectAttempt() {
	m.ReconnectAttempts.Add(1)
}

// RecordReconnectSuccess records a successful reconnection
func (m *Metrics) RecordReconnectSuccess() {
	m.ReconnectSuccesses.Add(1)
}

// RecordReconnectFailure records a failed reconnection
func (m *Metrics) RecordReconnectFailure() {
	m.ReconnectFailures.Add(1)
}

// RecordQueueEnqueue records a message being queued
func (m *Metrics) RecordQueueEnqueue() {
	m.QueueEnqueued.Add(1)
}

// RecordQueueSent records a queued message being sent
func (m *Metrics) RecordQueueSent() {
	m.QueueSent.Add(1)
}

// RecordQueueDropped records a queued message being dropped
func (m *Metrics) RecordQueueDropped() {
	m.QueueDropped.Add(1)
}

// RecordCompression records a compression operation
func (m *Metrics) RecordCompression(originalSize, compressedSize int) {
	m.CompressedMessages.Add(1)
	if originalSize > compressedSize {
		m.CompressionSaved.Add(int64(originalSize - compressedSize))
	}
}

// RecordDecompression records a decompression operation
func (m *Metrics) RecordDecompression() {
	m.DecompressedMessages.Add(1)
}

// DefaultMetrics is the default metrics instance
var DefaultMetrics = &Metrics{}
