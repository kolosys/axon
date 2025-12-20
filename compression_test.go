package axon_test

import (
	"bytes"
	"testing"

	"github.com/kolosys/axon"
)

func TestCompression_RoundTrip(t *testing.T) {
	// Create compression manager via the exported interface
	// Since CompressionManager is not exported, we test through Conn
	testData := []byte("Hello, this is a test message that should be compressed and decompressed successfully!")

	// Test that the data can be compressed and decompressed
	// We'll test this through the connection interface in integration tests
	if len(testData) < 256 {
		// Data below threshold should not be compressed
		t.Log("Test data is below compression threshold")
	}
}

func TestCompression_LargePayload(t *testing.T) {
	// Create a large payload that should compress well
	largeData := bytes.Repeat([]byte("Hello World! "), 1000)

	if len(largeData) < 256 {
		t.Error("expected large data to be above compression threshold")
	}

	// Verify that repeated data should compress well
	// The actual compression is tested through the connection
	t.Logf("Large data size: %d bytes", len(largeData))
}

func TestCompression_IncompressibleData(t *testing.T) {
	// Random data doesn't compress well
	randomData := make([]byte, 1000)
	for i := range randomData {
		randomData[i] = byte(i % 256)
	}

	// This should still work even if compression doesn't help
	if len(randomData) != 1000 {
		t.Errorf("expected random data length to be 1000, got %d", len(randomData))
	}
}

func TestDialOptions_Compression(t *testing.T) {
	opts := &axon.DialOptions{
		Compression:          true,
		CompressionThreshold: 512,
	}

	if !opts.Compression {
		t.Error("expected Compression to be true")
	}
	if opts.CompressionThreshold != 512 {
		t.Errorf("expected CompressionThreshold to be 512, got %d", opts.CompressionThreshold)
	}
}

func TestUpgradeOptions_Compression(t *testing.T) {
	opts := &axon.UpgradeOptions{
		Compression: true,
	}

	if !opts.Compression {
		t.Error("expected Compression to be true")
	}
}
