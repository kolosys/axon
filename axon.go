package axon

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	websocketGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

// Upgrader handles WebSocket connection upgrades
type Upgrader struct {
	readBufferSize    int
	writeBufferSize   int
	maxFrameSize      int
	maxMessageSize    int
	readDeadline      time.Duration
	writeDeadline     time.Duration
	pingInterval      time.Duration
	pongTimeout       time.Duration
	checkOrigin       func(r *http.Request) bool
	subprotocols      []string
	enableCompression bool
}

// NewUpgrader creates a new Upgrader with default settings
func NewUpgrader(opts *UpgradeOptions) *Upgrader {
	u := &Upgrader{
		readBufferSize:  4096,
		writeBufferSize: 4096,
		maxFrameSize:    4096,
		maxMessageSize:  1048576, // 1MB
	}

	if opts != nil {
		if opts.ReadBufferSize > 0 {
			u.readBufferSize = opts.ReadBufferSize
		}
		if opts.WriteBufferSize > 0 {
			u.writeBufferSize = opts.WriteBufferSize
		}
		if opts.MaxFrameSize > 0 {
			u.maxFrameSize = opts.MaxFrameSize
		}
		if opts.MaxMessageSize > 0 {
			u.maxMessageSize = opts.MaxMessageSize
		}
		u.readDeadline = opts.ReadDeadline
		u.writeDeadline = opts.WriteDeadline
		u.pingInterval = opts.PingInterval
		u.pongTimeout = opts.PongTimeout
		u.checkOrigin = opts.CheckOrigin
		u.subprotocols = opts.Subprotocols
		u.enableCompression = opts.Compression
	}

	return u
}

// Upgrade upgrades an HTTP connection to a WebSocket connection
func Upgrade[T any](w http.ResponseWriter, r *http.Request, opts *UpgradeOptions) (*Conn[T], error) {
	u := NewUpgrader(opts)
	return upgrade[T](u, w, r)
}

// upgrade performs the actual upgrade logic
func upgrade[T any](u *Upgrader, w http.ResponseWriter, r *http.Request) (*Conn[T], error) {
	if r.Method != http.MethodGet {
		return nil, ErrUpgradeRequired
	}

	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return nil, ErrUpgradeRequired
	}

	connection := r.Header.Get("Connection")
	if !strings.Contains(strings.ToLower(connection), "upgrade") {
		return nil, ErrUpgradeRequired
	}

	version := r.Header.Get("Sec-WebSocket-Version")
	if version != "13" {
		w.Header().Set("Sec-WebSocket-Version", "13")
		return nil, ErrInvalidHandshake
	}

	if u.checkOrigin != nil && !u.checkOrigin(r) {
		return nil, ErrInvalidOrigin
	}

	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return nil, ErrInvalidHandshake
	}

	requestedSubprotocol := r.Header.Get("Sec-WebSocket-Protocol")
	selectedSubprotocol := ""
	if requestedSubprotocol != "" && len(u.subprotocols) > 0 {
		requested := strings.Split(requestedSubprotocol, ",")
		for _, req := range requested {
			req = strings.TrimSpace(req)
			for _, supported := range u.subprotocols {
				if req == supported {
					selectedSubprotocol = supported
					break
				}
			}
			if selectedSubprotocol != "" {
				break
			}
		}
		if selectedSubprotocol == "" {
			return nil, ErrInvalidSubprotocol
		}
	}

	acceptKey := computeAcceptKey(key)

	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, ErrInvalidHandshake
	}

	conn, bufw, err := hj.Hijack()
	if err != nil {
		return nil, fmt.Errorf("axon: failed to hijack connection: %w", err)
	}

	response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\n"+
		"Upgrade: websocket\r\n"+
		"Connection: Upgrade\r\n"+
		"Sec-WebSocket-Accept: %s\r\n", acceptKey)

	if selectedSubprotocol != "" {
		response += fmt.Sprintf("Sec-WebSocket-Protocol: %s\r\n", selectedSubprotocol)
	}

	response += "\r\n"

	if _, err := bufw.WriteString(response); err != nil {
		conn.Close()
		return nil, fmt.Errorf("axon: failed to write response: %w", err)
	}

	if err := bufw.Flush(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("axon: failed to flush response: %w", err)
	}

	readBuf := getBuffer()
	writeBuf := getBuffer()
	reader := getReader(conn)
	writer := getWriter(conn)

	wsConn := &Conn[T]{
		conn:          conn,
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

	return wsConn, nil
}

// computeAcceptKey computes the WebSocket accept key (RFC 6455 Section 4.2.2)
func computeAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	h.Write([]byte(websocketGUID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
