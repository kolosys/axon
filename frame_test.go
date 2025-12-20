package axon_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/kolosys/axon"
)

func TestReadFrameHeader(t *testing.T) {
	// Simple text frame: FIN=1, opcode=1 (text), masked=0, payload_len=5
	frameData := []byte{0x81, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello"
	buf := make([]byte, 4096)

	frame, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read frame: %v", err)
	}

	if !frame.Fin {
		t.Error("FIN bit should be set")
	}
	if frame.Opcode != 0x1 {
		t.Errorf("expected opcode 0x1, got 0x%x", frame.Opcode)
	}
	if string(frame.Payload) != "Hello" {
		t.Errorf("expected payload 'Hello', got %q", string(frame.Payload))
	}
}

func TestReadFrameMasked(t *testing.T) {
	// Masked text frame
	mask := []byte{0x37, 0xfa, 0x21, 0x3d}
	payload := []byte("Hello")
	masked := make([]byte, len(payload))
	for i := range payload {
		masked[i] = payload[i] ^ mask[i%4]
	}

	frameData := make([]byte, 2+4+len(masked))
	frameData[0] = 0x81 // FIN=1, opcode=1
	frameData[1] = 0x85 // masked=1, payload_len=5
	copy(frameData[2:6], mask)
	copy(frameData[6:], masked)

	buf := make([]byte, 4096)
	frame, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read masked frame: %v", err)
	}

	if !frame.Masked {
		t.Error("frame should be marked as masked")
	}
	if string(frame.Payload) != "Hello" {
		t.Errorf("expected unmasked payload 'Hello', got %q", string(frame.Payload))
	}
}

func TestReadFrameExtendedLength(t *testing.T) {
	// Frame with 16-bit length (126)
	payload := make([]byte, 200)
	for i := range payload {
		payload[i] = byte(i % 256)
	}

	frameData := make([]byte, 2+2+len(payload))
	frameData[0] = 0x81 // FIN=1, opcode=1
	frameData[1] = 0x7E // masked=0, extended length (126)
	binary.BigEndian.PutUint16(frameData[2:4], uint16(len(payload)))
	copy(frameData[4:], payload)

	buf := make([]byte, 4096)
	frame, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read extended length frame: %v", err)
	}

	if len(frame.Payload) != len(payload) {
		t.Errorf("expected payload length %d, got %d", len(payload), len(frame.Payload))
	}
}

func TestReadFrame64BitLength(t *testing.T) {
	// Frame with 64-bit length (127)
	payload := make([]byte, 1000)
	for i := range payload {
		payload[i] = byte(i % 256)
	}

	frameData := make([]byte, 2+8+len(payload))
	frameData[0] = 0x81 // FIN=1, opcode=1
	frameData[1] = 0x7F // masked=0, extended length (127)
	// First 4 bytes are 0 (we only support 32-bit)
	binary.BigEndian.PutUint32(frameData[6:10], uint32(len(payload)))
	copy(frameData[10:], payload)

	buf := make([]byte, 4096)
	frame, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read 64-bit length frame: %v", err)
	}

	if len(frame.Payload) != len(payload) {
		t.Errorf("expected payload length %d, got %d", len(payload), len(frame.Payload))
	}
}

func TestReadFrameTooLarge(t *testing.T) {
	payload := make([]byte, 5000)
	frameData := make([]byte, 2+2+len(payload))
	frameData[0] = 0x81
	frameData[1] = 0x7E
	binary.BigEndian.PutUint16(frameData[2:4], uint16(len(payload)))
	copy(frameData[4:], payload)

	buf := make([]byte, 4096)
	_, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != axon.ErrFrameTooLarge {
		t.Errorf("expected ErrFrameTooLarge, got %v", err)
	}
}

func TestReadFrameInvalidRSV(t *testing.T) {
	// Frame with RSV1 bit set (invalid without extension)
	frameData := []byte{0xC1, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f} // RSV1=1
	buf := make([]byte, 4096)

	_, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != axon.ErrInvalidFrame {
		t.Errorf("expected ErrInvalidFrame, got %v", err)
	}
}

func TestReadFrameInvalidOpcode(t *testing.T) {
	// Frame with invalid opcode (0xB)
	frameData := []byte{0x8B, 0x05} // opcode=0xB
	buf := make([]byte, 4096)

	_, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != axon.ErrUnsupportedFrameType {
		t.Errorf("expected ErrUnsupportedFrameType, got %v", err)
	}
}

func TestReadFrameFragmentedControl(t *testing.T) {
	// Control frame with FIN=0 (invalid)
	frameData := []byte{0x08, 0x05} // FIN=0, opcode=8 (close)
	buf := make([]byte, 4096)

	_, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != axon.ErrFragmentedControlFrame {
		t.Errorf("expected ErrFragmentedControlFrame, got %v", err)
	}
}

func TestReadFrameIncomplete(t *testing.T) {
	// Incomplete frame
	frameData := []byte{0x81, 0x05, 0x48} // Missing bytes
	buf := make([]byte, 4096)

	_, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err == nil {
		t.Error("expected error for incomplete frame")
	}
}

func TestWriteFrame(t *testing.T) {
	frame := &axon.Frame{
		Fin:     true,
		Opcode:  0x1,
		Payload: []byte("Hello"),
	}

	var buf bytes.Buffer
	writeBuf := make([]byte, 4096)

	if err := axon.WriteFrame(&buf, writeBuf, frame); err != nil {
		t.Fatalf("failed to write frame: %v", err)
	}

	// Verify frame structure
	result := buf.Bytes()
	if len(result) < 7 {
		t.Fatalf("frame too short: %d bytes", len(result))
	}

	if result[0] != 0x81 {
		t.Errorf("expected first byte 0x81, got 0x%x", result[0])
	}
	if result[1] != 0x05 {
		t.Errorf("expected second byte 0x05, got 0x%x", result[1])
	}
	if string(result[2:]) != "Hello" {
		t.Errorf("expected payload 'Hello', got %q", string(result[2:]))
	}
}

func TestWriteFrameExtendedLength(t *testing.T) {
	payload := make([]byte, 200)
	for i := range payload {
		payload[i] = byte(i % 256)
	}

	frame := &axon.Frame{
		Fin:     true,
		Opcode:  0x2,
		Payload: payload,
	}

	var buf bytes.Buffer
	writeBuf := make([]byte, 4096)

	if err := axon.WriteFrame(&buf, writeBuf, frame); err != nil {
		t.Fatalf("failed to write extended length frame: %v", err)
	}

	result := buf.Bytes()
	if len(result) != 4+len(payload) {
		t.Errorf("expected %d bytes, got %d", 4+len(payload), len(result))
	}

	// Verify extended length encoding
	if result[1] != 0x7E {
		t.Errorf("expected extended length marker 0x7E, got 0x%x", result[1])
	}
}

func TestWriteFrameMasked(t *testing.T) {
	mask := []byte{0x37, 0xfa, 0x21, 0x3d}
	frame := &axon.Frame{
		Fin:     true,
		Opcode:  0x1,
		Masked:  true,
		MaskKey: mask,
		Payload: []byte("Hello"),
	}

	var buf bytes.Buffer
	writeBuf := make([]byte, 4096)

	if err := axon.WriteFrame(&buf, writeBuf, frame); err != nil {
		t.Fatalf("failed to write masked frame: %v", err)
	}

	result := buf.Bytes()
	// Should have header + mask + payload
	if len(result) != 2+4+5 {
		t.Errorf("expected %d bytes, got %d", 2+4+5, len(result))
	}

	// Verify mask key is present
	if !bytes.Equal(result[2:6], mask) {
		t.Error("mask key not found in frame")
	}
}

func TestWriteFrameRSV(t *testing.T) {
	frame := &axon.Frame{
		Fin:     true,
		Rsv1:    true,
		Rsv2:    true,
		Rsv3:    true,
		Opcode:  0x1,
		Payload: []byte("test"),
	}

	var buf bytes.Buffer
	writeBuf := make([]byte, 4096)

	if err := axon.WriteFrame(&buf, writeBuf, frame); err != nil {
		t.Fatalf("failed to write frame with RSV bits: %v", err)
	}

	result := buf.Bytes()
	// RSV bits should be set
	if (result[0] & 0x70) != 0x70 {
		t.Errorf("expected RSV bits 0x70, got 0x%x", result[0]&0x70)
	}
}

func TestWriteFrameLargePayload(t *testing.T) {
	payload := make([]byte, 100000)
	frame := &axon.Frame{
		Fin:     true,
		Opcode:  0x2,
		Payload: payload,
	}

	var buf bytes.Buffer
	writeBuf := make([]byte, 4096)

	if err := axon.WriteFrame(&buf, writeBuf, frame); err != nil {
		t.Fatalf("failed to write large frame: %v", err)
	}

	result := buf.Bytes()
	// Should use 64-bit length encoding
	if result[1] != 0x7F {
		t.Errorf("expected 64-bit length marker 0x7F, got 0x%x", result[1])
	}
}

func TestWriteFrameWriteError(t *testing.T) {
	frame := &axon.Frame{
		Fin:     true,
		Opcode:  0x1,
		Payload: []byte("test"),
	}

	// Create a writer that fails
	errWriter := &errorWriter{}
	writeBuf := make([]byte, 4096)

	err := axon.WriteFrame(errWriter, writeBuf, frame)
	if err == nil {
		t.Error("expected error from failing writer")
	}
}

type errorWriter struct{}

func (e *errorWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestReadFrameContinuation(t *testing.T) {
	// First frame: FIN=0, opcode=1 (text), payload="Hel"
	frame1 := []byte{0x01, 0x03, 0x48, 0x65, 0x6c}
	// Continuation frame: FIN=1, opcode=0, payload="lo"
	frame2 := []byte{0x80, 0x02, 0x6c, 0x6f}

	combined := append(frame1, frame2...)
	buf := make([]byte, 4096)

	// Read first frame
	frame, err := axon.ReadFrame(bytes.NewReader(combined), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read first frame: %v", err)
	}

	if frame.Fin {
		t.Error("first frame should not have FIN bit set")
	}
	if frame.Opcode != 0x1 {
		t.Errorf("expected opcode 0x1, got 0x%x", frame.Opcode)
	}
	if string(frame.Payload) != "Hel" {
		t.Errorf("expected 'Hel', got %q", string(frame.Payload))
	}

	// Read continuation frame
	frame, err = axon.ReadFrame(bytes.NewReader(combined[len(frame1):]), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read continuation frame: %v", err)
	}

	if !frame.Fin {
		t.Error("continuation frame should have FIN bit set")
	}
	if frame.Opcode != 0x0 {
		t.Errorf("expected opcode 0x0 (continuation), got 0x%x", frame.Opcode)
	}
	if string(frame.Payload) != "lo" {
		t.Errorf("expected 'lo', got %q", string(frame.Payload))
	}
}

func TestReadFrameContinuationIsolated(t *testing.T) {
	// Continuation frame can be read in isolation (validation happens at connection level)
	frameData := []byte{0x80, 0x02, 0x6c, 0x6f} // FIN=1, opcode=0 (continuation)
	buf := make([]byte, 4096)

	frame, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read continuation frame: %v", err)
	}

	if frame.Opcode != 0x0 {
		t.Errorf("expected opcode 0x0 (continuation), got 0x%x", frame.Opcode)
	}
	if !frame.Fin {
		t.Error("continuation frame should have FIN bit set")
	}
}

func TestWriteFrameEmptyPayload(t *testing.T) {
	frame := &axon.Frame{
		Fin:     true,
		Opcode:  0x1,
		Payload: []byte{},
	}

	var buf bytes.Buffer
	writeBuf := make([]byte, 4096)

	if err := axon.WriteFrame(&buf, writeBuf, frame); err != nil {
		t.Fatalf("failed to write empty frame: %v", err)
	}

	result := buf.Bytes()
	if len(result) < 2 {
		t.Fatalf("frame too short: %d bytes", len(result))
	}

	if result[0] != 0x81 {
		t.Errorf("expected first byte 0x81, got 0x%x", result[0])
	}
	if result[1] != 0x00 {
		t.Errorf("expected second byte 0x00 (empty payload), got 0x%x", result[1])
	}
}

func TestReadFramePingPong(t *testing.T) {
	// Ping frame
	pingData := []byte{0x89, 0x05, 0x70, 0x69, 0x6e, 0x67, 0x20} // "ping "
	buf := make([]byte, 4096)

	frame, err := axon.ReadFrame(bytes.NewReader(pingData), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read ping frame: %v", err)
	}

	if frame.Opcode != 0x9 {
		t.Errorf("expected opcode 0x9 (ping), got 0x%x", frame.Opcode)
	}
	if !frame.Fin {
		t.Error("ping frame should have FIN bit set")
	}

	// Pong frame
	pongData := []byte{0x8A, 0x05, 0x70, 0x69, 0x6e, 0x67, 0x20} // "ping "
	frame, err = axon.ReadFrame(bytes.NewReader(pongData), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read pong frame: %v", err)
	}

	if frame.Opcode != 0xA {
		t.Errorf("expected opcode 0xA (pong), got 0x%x", frame.Opcode)
	}
}

func TestReadFrameClose(t *testing.T) {
	// Close frame with code 1000 and reason "Normal closure"
	closeData := []byte{0x88, 0x0D, 0x03, 0xE8, 0x4E, 0x6F, 0x72, 0x6D, 0x61, 0x6C, 0x20, 0x63, 0x6C, 0x6F, 0x73, 0x75, 0x72, 0x65}
	buf := make([]byte, 4096)

	frame, err := axon.ReadFrame(bytes.NewReader(closeData), buf, 4096)
	if err != nil {
		t.Fatalf("failed to read close frame: %v", err)
	}

	if frame.Opcode != 0x8 {
		t.Errorf("expected opcode 0x8 (close), got 0x%x", frame.Opcode)
	}
	if len(frame.Payload) < 2 {
		t.Error("close frame should have at least 2 bytes (code)")
	}
}

func TestWriteFrameClose(t *testing.T) {
	closePayload := make([]byte, 2+13)
	binary.BigEndian.PutUint16(closePayload[:2], 1000)
	copy(closePayload[2:], []byte("Normal closure"))

	frame := &axon.Frame{
		Fin:     true,
		Opcode:  0x8,
		Payload: closePayload,
	}

	var buf bytes.Buffer
	writeBuf := make([]byte, 4096)

	if err := axon.WriteFrame(&buf, writeBuf, frame); err != nil {
		t.Fatalf("failed to write close frame: %v", err)
	}

	result := buf.Bytes()
	if result[0] != 0x88 {
		t.Errorf("expected first byte 0x88, got 0x%x", result[0])
	}
}

func TestReadFrame64BitLengthTooLarge(t *testing.T) {
	// Frame with 64-bit length where first 4 bytes are non-zero (too large)
	frameData := make([]byte, 2+8)
	frameData[0] = 0x81
	frameData[1] = 0x7F // 64-bit length
	frameData[2] = 0x01 // Non-zero - should trigger ErrFrameTooLarge
	buf := make([]byte, 4096)

	_, err := axon.ReadFrame(bytes.NewReader(frameData), buf, 4096)
	if err != axon.ErrFrameTooLarge {
		t.Errorf("expected ErrFrameTooLarge, got %v", err)
	}
}
