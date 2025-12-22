package axon

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Conn represents a WebSocket connection with type-safe message handling
type Conn[T any] struct {
	conn          net.Conn
	reader        *bufio.Reader
	writer        *bufio.Writer
	readBuf       []byte
	writeBuf      []byte
	upgrader      *Upgrader
	closed        int32
	closeOnce     sync.Once
	closeCode     int
	closeReason   string
	readDeadline  time.Duration
	writeDeadline time.Duration
	pingInterval  time.Duration
	pongTimeout   time.Duration
	pingTicker    *time.Ticker
	pingStop      chan struct{}
	pingWg        sync.WaitGroup
	isClient      bool
	compression   *CompressionManager
	writeMu       sync.Mutex
}

// Read reads a complete message from the connection
func (c *Conn[T]) Read(ctx context.Context) (T, error) {
	var zero T

	if atomic.LoadInt32(&c.closed) != 0 {
		return zero, ErrConnectionClosed
	}

	deadline := c.readDeadline
	if deadline == 0 {
		deadline = 30 * time.Second // Default timeout
	}

	if ctx != nil {
		if ctxDeadline, ok := ctx.Deadline(); ok {
			deadline = time.Until(ctxDeadline)
		}
		if ctx.Err() != nil {
			return zero, ErrContextCanceled
		}
	}

	if err := c.conn.SetReadDeadline(time.Now().Add(deadline)); err != nil {
		return zero, err
	}

	var messagePayload []byte
	firstFrame := true

	for {
		frame, err := readFrame(c.reader, c.readBuf, c.upgrader.maxFrameSize)
		if err != nil {
			if err == io.EOF {
				return zero, ErrConnectionClosed
			}
			return zero, err
		}

		switch frame.Opcode {
		case opContinuation:
			if firstFrame {
				return zero, ErrInvalidFrame
			}
		case opClose:
			code := 1000 // Normal closure
			reason := ""
			if len(frame.Payload) >= 2 {
				code = int(binary.BigEndian.Uint16(frame.Payload[:2]))
				if len(frame.Payload) > 2 {
					reason = string(frame.Payload[2:])
				}
			}
			c.closeOnce.Do(func() {
				atomic.StoreInt32(&c.closed, 1)
				c.closeCode = code
				c.closeReason = reason
			})
			return zero, ErrConnectionClosed

		case opPing:
			pongFrame := &Frame{
				Fin:     true,
				Opcode:  opPong,
				Payload: frame.Payload,
			}
			if err := writeFrame(c.writer, c.writeBuf, pongFrame); err != nil {
				return zero, err
			}
			if err := c.writer.Flush(); err != nil {
				return zero, err
			}
			continue

		case opPong:
			continue
		case opText, opBinary:
			if !firstFrame {
				return zero, ErrInvalidFrame
			}
			firstFrame = false
		default:
			return zero, ErrUnsupportedFrameType
		}

		messagePayload = append(messagePayload, frame.Payload...)

		if len(messagePayload) > c.upgrader.maxMessageSize {
			return zero, ErrMessageTooLarge
		}

		if frame.Fin {
			break
		}
	}

	var msg T
	if len(messagePayload) == 0 {
		return zero, nil
	}

	// Decompress if compression is enabled and message was compressed
	if c.compression != nil && c.compression.enabled {
		decompressed, err := c.compression.Decompress(messagePayload)
		if err != nil {
			return zero, err
		}
		messagePayload = decompressed
	}

	if err := json.Unmarshal(messagePayload, &msg); err != nil {
		switch v := any(&msg).(type) {
		case *[]byte:
			*v = messagePayload
			return msg, nil
		case *string:
			*v = string(messagePayload)
			return msg, nil
		}
		return zero, ErrDeserializationFailed
	}

	return msg, nil
}

// IsClosed returns true if the connection has been closed
func (c *Conn[T]) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) != 0
}

// CloseCode returns the close code if the connection was closed
func (c *Conn[T]) CloseCode() int {
	return c.closeCode
}

// CloseReason returns the close reason if the connection was closed
func (c *Conn[T]) CloseReason() string {
	return c.closeReason
}

// Write writes a message to the connection
func (c *Conn[T]) Write(ctx context.Context, msg T) error {
	if atomic.LoadInt32(&c.closed) != 0 {
		return ErrConnectionClosed
	}

	deadline := c.writeDeadline
	if deadline == 0 {
		deadline = 30 * time.Second // Default timeout
	}

	if ctx != nil {
		if ctxDeadline, ok := ctx.Deadline(); ok {
			deadline = time.Until(ctxDeadline)
		}
		if ctx.Err() != nil {
			return ErrContextCanceled
		}
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if err := c.conn.SetWriteDeadline(time.Now().Add(deadline)); err != nil {
		return err
	}

	var payload []byte
	var err error

	if payload, err = json.Marshal(msg); err != nil {
		switch v := any(msg).(type) {
		case []byte:
			payload = v
		case string:
			payload = []byte(v)
		default:
			return ErrSerializationFailed
		}
	}

	if len(payload) > c.upgrader.maxMessageSize {
		return ErrMessageTooLarge
	}

	var opcode byte
	switch any(msg).(type) {
	case string:
		opcode = opText
	case []byte:
		opcode = opBinary
	default:
		opcode = opText // JSON is text frame
	}

	// Compress if compression is enabled and payload is large enough
	compressed := false
	if c.compression != nil && c.compression.ShouldCompress(len(payload)) {
		compressedPayload, err := c.compression.Compress(payload)
		if err == nil && len(compressedPayload) < len(payload) {
			payload = compressedPayload
			compressed = true
		}
	}

	frame := &Frame{
		Fin:     true,
		Rsv1:    compressed, // RSV1 indicates compression
		Opcode:  opcode,
		Masked:  c.isClient, // Clients must mask frames
		Payload: payload,
	}

	// Generate mask key for client connections
	if c.isClient {
		frame.MaskKey = make([]byte, 4)
		if _, err := rand.Read(frame.MaskKey); err != nil {
			return fmt.Errorf("axon: failed to generate mask key: %w", err)
		}
		// Mask the payload
		maskedPayload := make([]byte, len(payload))
		copy(maskedPayload, payload)
		maskBytes(maskedPayload, frame.MaskKey)
		frame.Payload = maskedPayload
	}

	if err := writeFrame(c.writer, c.writeBuf, frame); err != nil {
		return err
	}

	return c.writer.Flush()
}

// Close closes the connection with the given code and reason
func (c *Conn[T]) Close(code int, reason string) error {
	var closeErr error

	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closed, 1)
		c.closeCode = code
		c.closeReason = reason

		if c.pingStop != nil {
			close(c.pingStop)
			c.pingWg.Wait()
			if c.pingTicker != nil {
				c.pingTicker.Stop()
			}
		}

		closePayload := make([]byte, 2+len(reason))
		binary.BigEndian.PutUint16(closePayload[:2], uint16(code))
		copy(closePayload[2:], reason)

		closeFrame := &Frame{
			Fin:     true,
			Opcode:  opClose,
			Payload: closePayload,
		}

		// Set a short deadline to avoid blocking on close frame write
		c.conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))

		if err := writeFrame(c.writer, c.writeBuf, closeFrame); err != nil {
			// Ignore write errors on close - connection may already be dead
			c.conn.Close()
			closeErr = nil
			return
		}

		if err := c.writer.Flush(); err != nil {
			// Ignore flush errors on close - connection may already be dead
			c.conn.Close()
			closeErr = nil
			return
		}

		closeErr = c.conn.Close()

		putBuffer(c.readBuf)
		putBuffer(c.writeBuf)
		putReader(c.reader)
		putWriter(c.writer)
	})

	return closeErr
}

// startPingLoop starts the ping/pong keepalive loop
func (c *Conn[T]) startPingLoop() {
	if c.pingInterval == 0 {
		return
	}

	c.pingStop = make(chan struct{})
	c.pingTicker = time.NewTicker(c.pingInterval)

	c.pingWg.Add(1)
	go func() {
		defer c.pingWg.Done()

		for {
			select {
			case <-c.pingTicker.C:
				pingFrame := &Frame{
					Fin:     true,
					Opcode:  opPing,
					Payload: []byte("ping"),
				}

				if err := c.conn.SetWriteDeadline(time.Now().Add(c.pongTimeout)); err == nil {
					writeFrame(c.writer, c.writeBuf, pingFrame)
					c.writer.Flush()
				}

			case <-c.pingStop:
				return
			}
		}
	}()
}
