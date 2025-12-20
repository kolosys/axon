package axon

import (
	"bytes"
	"compress/flate"
	"io"
	"sync"
)

// CompressionManager handles per-message compression (RFC 7692)
type CompressionManager struct {
	enabled   bool
	threshold int // Minimum size to compress

	// Compressor resources
	compressorMu sync.Mutex
	compressor   *flate.Writer
	compressBuf  bytes.Buffer

	// Decompressor resources
	decompressorMu sync.Mutex
	decompressor   io.ReadCloser
	decompressBuf  bytes.Buffer
}

// newCompressionManager creates a new CompressionManager
func newCompressionManager(threshold int) *CompressionManager {
	if threshold <= 0 {
		threshold = 256
	}
	return &CompressionManager{
		enabled:   true,
		threshold: threshold,
	}
}

// ShouldCompress returns true if the payload should be compressed
func (cm *CompressionManager) ShouldCompress(payloadSize int) bool {
	return cm.enabled && payloadSize >= cm.threshold
}

// Compress compresses the payload using DEFLATE
func (cm *CompressionManager) Compress(data []byte) ([]byte, error) {
	cm.compressorMu.Lock()
	defer cm.compressorMu.Unlock()

	cm.compressBuf.Reset()

	// Create compressor if not exists
	if cm.compressor == nil {
		var err error
		cm.compressor, err = flate.NewWriter(&cm.compressBuf, flate.BestSpeed)
		if err != nil {
			return nil, err
		}
	} else {
		cm.compressor.Reset(&cm.compressBuf)
	}

	// Write data to compressor
	if _, err := cm.compressor.Write(data); err != nil {
		return nil, err
	}

	// Flush the compressor
	if err := cm.compressor.Flush(); err != nil {
		return nil, err
	}

	// Get compressed data
	compressed := cm.compressBuf.Bytes()

	// Per RFC 7692, remove the trailing 0x00 0x00 0xff 0xff
	if len(compressed) >= 4 {
		tail := compressed[len(compressed)-4:]
		if tail[0] == 0x00 && tail[1] == 0x00 && tail[2] == 0xff && tail[3] == 0xff {
			compressed = compressed[:len(compressed)-4]
		}
	}

	// Make a copy to avoid buffer reuse issues
	result := make([]byte, len(compressed))
	copy(result, compressed)

	return result, nil
}

// Decompress decompresses the payload using DEFLATE
func (cm *CompressionManager) Decompress(data []byte) ([]byte, error) {
	cm.decompressorMu.Lock()
	defer cm.decompressorMu.Unlock()

	// Per RFC 7692, append the trailing 0x00 0x00 0xff 0xff
	dataWithTail := make([]byte, len(data)+4)
	copy(dataWithTail, data)
	dataWithTail[len(data)] = 0x00
	dataWithTail[len(data)+1] = 0x00
	dataWithTail[len(data)+2] = 0xff
	dataWithTail[len(data)+3] = 0xff

	cm.decompressBuf.Reset()
	reader := bytes.NewReader(dataWithTail)

	// Create decompressor
	cm.decompressor = flate.NewReader(reader)
	defer cm.decompressor.Close()

	// Read decompressed data
	if _, err := io.Copy(&cm.decompressBuf, cm.decompressor); err != nil {
		return nil, err
	}

	// Make a copy to avoid buffer reuse issues
	result := make([]byte, cm.decompressBuf.Len())
	copy(result, cm.decompressBuf.Bytes())

	return result, nil
}

// Close releases compression resources
func (cm *CompressionManager) Close() error {
	cm.compressorMu.Lock()
	if cm.compressor != nil {
		cm.compressor.Close()
		cm.compressor = nil
	}
	cm.compressorMu.Unlock()

	cm.decompressorMu.Lock()
	if cm.decompressor != nil {
		cm.decompressor.Close()
		cm.decompressor = nil
	}
	cm.decompressorMu.Unlock()

	return nil
}
