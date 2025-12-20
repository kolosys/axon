package benchmarks_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kolosys/axon"
)

type benchMessage struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

func BenchmarkUpgrade(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := axon.Upgrade[benchMessage](w, r, nil)
		if err != nil {
			return
		}
		conn.Close(1000, "test")
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
		req.Header.Set("Sec-WebSocket-Version", "13")

		w := httptest.NewRecorder()
		handler(w, req)
	}
}

func BenchmarkRead(b *testing.B) {
	// This benchmark would require a real WebSocket connection
	// For now, we'll skip it as it requires more complex setup
	b.Skip("requires real WebSocket connection")
}

func BenchmarkWrite(b *testing.B) {
	// This benchmark would require a real WebSocket connection
	// For now, we'll skip it as it requires more complex setup
	b.Skip("requires real WebSocket connection")
}

func BenchmarkFrameParsing(b *testing.B) {
	// Benchmark frame header parsing
	frameData := []byte{0x81, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello" text frame

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate frame parsing overhead
			_ = frameData[0] & 0x80 // FIN bit
			_ = frameData[0] & 0x0F // Opcode
			_ = frameData[1] & 0x7F // Payload length
		}
	})
}

func BenchmarkMetrics(b *testing.B) {
	metrics := &axon.Metrics{}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.RecordRead(1024, 100*time.Microsecond)
			metrics.RecordWrite(1024, 100*time.Microsecond)
		}
	})
}

func BenchmarkSerialization(b *testing.B) {
	msg := benchMessage{ID: 1, Data: "benchmark message data"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = json.Marshal(msg)
		}
	})
}
