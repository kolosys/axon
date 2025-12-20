package axon_test

import (
	"errors"
	"testing"

	"github.com/kolosys/axon"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"InvalidHandshake", axon.ErrInvalidHandshake},
		{"UpgradeRequired", axon.ErrUpgradeRequired},
		{"InvalidOrigin", axon.ErrInvalidOrigin},
		{"InvalidSubprotocol", axon.ErrInvalidSubprotocol},
		{"ConnectionClosed", axon.ErrConnectionClosed},
		{"FrameTooLarge", axon.ErrFrameTooLarge},
		{"MessageTooLarge", axon.ErrMessageTooLarge},
		{"InvalidFrame", axon.ErrInvalidFrame},
		{"InvalidMask", axon.ErrInvalidMask},
		{"UnsupportedFrameType", axon.ErrUnsupportedFrameType},
		{"FragmentedControlFrame", axon.ErrFragmentedControlFrame},
		{"InvalidCloseCode", axon.ErrInvalidCloseCode},
		{"ReadDeadlineExceeded", axon.ErrReadDeadlineExceeded},
		{"WriteDeadlineExceeded", axon.ErrWriteDeadlineExceeded},
		{"ContextCanceled", axon.ErrContextCanceled},
		{"SerializationFailed", axon.ErrSerializationFailed},
		{"DeserializationFailed", axon.ErrDeserializationFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatal("error should not be nil")
			}
			if !errors.Is(tt.err, tt.err) {
				t.Fatal("error should be comparable")
			}
		})
	}
}
