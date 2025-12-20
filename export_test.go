package axon

import "net"

// Export internal functions for testing
var (
	ReadFrame  = readFrame
	WriteFrame = writeFrame
	GetBuffer  = getBuffer
	PutBuffer  = putBuffer
	GetReader  = getReader
	PutReader  = putReader
	GetWriter  = getWriter
	PutWriter  = putWriter
)

// NewTestConn creates a Conn for testing using net.Pipe
func NewTestConn[T any](opts *UpgradeOptions) (*Conn[T], net.Conn, error) {
	if opts == nil {
		opts = &UpgradeOptions{}
	}

	u := NewUpgrader(opts)
	clientConn, serverConn := net.Pipe()

	readBuf := getBuffer()
	writeBuf := getBuffer()
	reader := getReader(serverConn)
	writer := getWriter(serverConn)

	wsConn := &Conn[T]{
		conn:          serverConn,
		reader:        reader,
		writer:        writer,
		readBuf:       readBuf,
		writeBuf:      writeBuf,
		upgrader:      u,
		readDeadline:  u.readDeadline,
		writeDeadline: u.writeDeadline,
		pingInterval:  u.pingInterval,
		pongTimeout:   u.pongTimeout,
	}

	if u.pingInterval > 0 {
		wsConn.startPingLoop()
	}

	return wsConn, clientConn, nil
}
