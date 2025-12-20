package axon_test

import (
	"bytes"
	"testing"

	"github.com/kolosys/axon"
)

func TestBufferPool(t *testing.T) {
	// Get buffer from pool
	buf1 := axon.GetBuffer()
	if len(buf1) == 0 {
		t.Error("buffer should not be empty")
	}

	// Put it back
	axon.PutBuffer(buf1)

	// Get another buffer - should reuse if pool works
	buf2 := axon.GetBuffer()
	if len(buf2) == 0 {
		t.Error("buffer should not be empty")
	}

	// Put back
	axon.PutBuffer(buf2)
}

func TestBufferPoolSmallBuffer(t *testing.T) {
	// Small buffer should not be put back
	smallBuf := make([]byte, 100)
	axon.PutBuffer(smallBuf)
	// Should not panic
}

func TestReaderPool(t *testing.T) {
	r := bytes.NewReader([]byte("test"))
	br := axon.GetReader(r)

	if br == nil {
		t.Fatal("reader should not be nil")
	}

	// Read some data
	buf := make([]byte, 4)
	n, err := br.Read(buf)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if n != 4 {
		t.Errorf("expected 4 bytes, got %d", n)
	}

	// Put back
	axon.PutReader(br)
}

func TestWriterPool(t *testing.T) {
	var buf bytes.Buffer
	bw := axon.GetWriter(&buf)

	if bw == nil {
		t.Fatal("writer should not be nil")
	}

	// Write some data
	n, err := bw.Write([]byte("test"))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
	if n != 4 {
		t.Errorf("expected 4 bytes written, got %d", n)
	}

	// Flush
	if err := bw.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	if buf.String() != "test" {
		t.Errorf("expected 'test', got %q", buf.String())
	}

	// Put back
	axon.PutWriter(bw)
}
