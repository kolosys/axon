package axon

import (
	"encoding/binary"
	"io"
)

const (
	// Frame opcodes (RFC 6455 Section 5.2)
	opContinuation = 0x0
	opText         = 0x1
	opBinary       = 0x2
	opClose        = 0x8
	opPing         = 0x9
	opPong         = 0xA

	// Frame flags
	finMask = 0x80
	rsvMask = 0x70
	opMask  = 0x0F

	// Maximum frame header size (2 bytes base + 8 bytes extended length + 4 bytes mask)
	maxFrameHeaderSize = 14
)

// Frame represents a WebSocket frame
type Frame struct {
	Fin     bool
	Rsv1    bool
	Rsv2    bool
	Rsv3    bool
	Opcode  byte
	Masked  bool
	MaskKey []byte
	Payload []byte
}

// readFrameHeader reads and parses a WebSocket frame header without allocations
func readFrameHeader(r io.Reader, buf []byte) (*Frame, error) {
	if _, err := io.ReadFull(r, buf[:2]); err != nil {
		return nil, err
	}

	frame := &Frame{
		Fin:    (buf[0] & finMask) != 0,
		Rsv1:   (buf[0] & 0x40) != 0,
		Rsv2:   (buf[0] & 0x20) != 0,
		Rsv3:   (buf[0] & 0x10) != 0,
		Opcode: buf[0] & opMask,
		Masked: (buf[1] & 0x80) != 0,
	}

	// RSV bits must be 0 unless extension negotiated
	if (buf[0] & rsvMask) != 0 {
		return nil, ErrInvalidFrame
	}

	if frame.Opcode > 0x7 && frame.Opcode < 0x8 {
		return nil, ErrUnsupportedFrameType
	}
	if frame.Opcode > 0xA {
		return nil, ErrUnsupportedFrameType
	}

	if frame.Opcode >= 0x8 && !frame.Fin {
		return nil, ErrFragmentedControlFrame
	}

	payloadLen := int(buf[1] & 0x7F)
	headerSize := 2

	switch payloadLen {
	case 126:
		if _, err := io.ReadFull(r, buf[2:4]); err != nil {
			return nil, err
		}
		payloadLen = int(binary.BigEndian.Uint16(buf[2:4]))
		headerSize = 4
	case 127:
		if _, err := io.ReadFull(r, buf[2:10]); err != nil {
			return nil, err
		}
		if buf[2] != 0 || buf[3] != 0 || buf[4] != 0 || buf[5] != 0 {
			return nil, ErrFrameTooLarge
		}
		payloadLen = int(binary.BigEndian.Uint32(buf[6:10]))
		headerSize = 10
	}

	if frame.Masked {
		if _, err := io.ReadFull(r, buf[headerSize:headerSize+4]); err != nil {
			return nil, err
		}
		frame.MaskKey = buf[headerSize : headerSize+4]
		headerSize += 4
	}

	frame.Payload = make([]byte, payloadLen)
	return frame, nil
}

// readFrame reads a complete frame including payload
func readFrame(r io.Reader, buf []byte, maxSize int) (*Frame, error) {
	frame, err := readFrameHeader(r, buf)
	if err != nil {
		return nil, err
	}

	if len(frame.Payload) > maxSize {
		return nil, ErrFrameTooLarge
	}

	if len(frame.Payload) > 0 {
		if _, err := io.ReadFull(r, frame.Payload); err != nil {
			return nil, err
		}

		if frame.Masked {
			maskBytes(frame.Payload, frame.MaskKey)
		}
	}

	return frame, nil
}

// writeFrame writes a frame header and payload
func writeFrame(w io.Writer, buf []byte, frame *Frame) error {
	headerSize := 2
	buf[0] = 0

	if frame.Fin {
		buf[0] |= finMask
	}

	if frame.Rsv1 {
		buf[0] |= 0x40
	}
	if frame.Rsv2 {
		buf[0] |= 0x20
	}
	if frame.Rsv3 {
		buf[0] |= 0x10
	}

	buf[0] |= frame.Opcode & opMask

	payloadLen := len(frame.Payload)
	if frame.Masked {
		buf[1] = 0x80
	} else {
		buf[1] = 0
	}

	switch {
	case payloadLen < 126:
		buf[1] |= byte(payloadLen)
	case payloadLen < 65536:
		buf[1] |= 126
		binary.BigEndian.PutUint16(buf[2:4], uint16(payloadLen))
		headerSize = 4
	default:
		buf[1] |= 127
		binary.BigEndian.PutUint64(buf[2:10], uint64(payloadLen))
		headerSize = 10
	}

	if frame.Masked {
		copy(buf[headerSize:headerSize+4], frame.MaskKey)
		headerSize += 4
	}

	if _, err := w.Write(buf[:headerSize]); err != nil {
		return err
	}

	if len(frame.Payload) > 0 {
		if _, err := w.Write(frame.Payload); err != nil {
			return err
		}
	}

	return nil
}

// maskBytes applies XOR masking to payload (RFC 6455 Section 5.3)
func maskBytes(payload []byte, mask []byte) {
	for i := range payload {
		payload[i] ^= mask[i%4]
	}
}
