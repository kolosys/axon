package axon

import (
	"bufio"
	"sync"
)

// bufferPool manages reusable buffers for zero-allocation operations
// Buffer size must accommodate max frame header (14 bytes) plus reasonable payload
var bufferPool = sync.Pool{
	New: func() any {
		buf := make([]byte, maxFrameHeaderSize+4096)
		return &buf
	},
}

// getBuffer retrieves a buffer from the pool
func getBuffer() []byte {
	return *bufferPool.Get().(*[]byte)
}

// putBuffer returns a buffer to the pool
func putBuffer(buf []byte) {
	if cap(buf) >= 4096 {
		// Reset length to capacity to avoid retaining references
		buf = buf[:cap(buf)]
		bufferPool.Put(&buf)
	}
}

// readerPool manages reusable bufio.Reader instances
var readerPool = sync.Pool{
	New: func() any {
		return bufio.NewReaderSize(nil, 4096)
	},
}

// getReader retrieves a reader from the pool and resets it with the given reader
func getReader(r interface{ Read([]byte) (int, error) }) *bufio.Reader {
	br := readerPool.Get().(*bufio.Reader)
	br.Reset(r)
	return br
}

// putReader returns a reader to the pool
func putReader(br *bufio.Reader) {
	br.Reset(nil)
	readerPool.Put(br)
}

// writerPool manages reusable bufio.Writer instances
var writerPool = sync.Pool{
	New: func() any {
		return bufio.NewWriterSize(nil, 4096)
	},
}

// getWriter retrieves a writer from the pool and resets it with the given writer
func getWriter(w interface{ Write([]byte) (int, error) }) *bufio.Writer {
	bw := writerPool.Get().(*bufio.Writer)
	bw.Reset(w)
	return bw
}

// putWriter returns a writer to the pool
func putWriter(bw *bufio.Writer) {
	bw.Reset(nil)
	writerPool.Put(bw)
}
