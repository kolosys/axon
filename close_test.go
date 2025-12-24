package axon_test

import (
	"testing"

	"github.com/kolosys/axon"
)

func TestCloseCodeString(t *testing.T) {
	tests := []struct {
		code     axon.CloseCode
		expected string
	}{
		{axon.CloseNormalClosure, "Normal Closure"},
		{axon.CloseGoingAway, "Going Away"},
		{axon.CloseProtocolError, "Protocol Error"},
		{axon.CloseUnsupportedData, "Unsupported Data"},
		{axon.CloseInternalError, "Internal Error"},
		{axon.CloseServiceRestart, "Service Restart"},
		{axon.CloseCode(3000), "Library-Defined"},
		{axon.CloseCode(4014), "Application-Defined"},
		{axon.CloseCode(9999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.code.String(); got != tt.expected {
				t.Errorf("CloseCode(%d).String() = %q, want %q", tt.code, got, tt.expected)
			}
		})
	}
}

func TestCloseCodeIsReserved(t *testing.T) {
	reserved := []axon.CloseCode{
		axon.CloseNoStatusReceived,
		axon.CloseAbnormalClosure,
		axon.CloseTLSHandshake,
	}

	notReserved := []axon.CloseCode{
		axon.CloseNormalClosure,
		axon.CloseGoingAway,
		axon.CloseProtocolError,
		axon.CloseInternalError,
	}

	for _, code := range reserved {
		if !code.IsReserved() {
			t.Errorf("CloseCode(%d).IsReserved() = false, want true", code)
		}
	}

	for _, code := range notReserved {
		if code.IsReserved() {
			t.Errorf("CloseCode(%d).IsReserved() = true, want false", code)
		}
	}
}

func TestCloseCodeIsValid(t *testing.T) {
	valid := []axon.CloseCode{
		axon.CloseNormalClosure,
		axon.CloseGoingAway,
		axon.CloseProtocolError,
		axon.CloseCode(3000), // Library-defined
		axon.CloseCode(4000), // Application-defined
		axon.CloseCode(4999), // Application-defined
	}

	invalid := []axon.CloseCode{
		axon.CloseNoStatusReceived, // Reserved
		axon.CloseAbnormalClosure,  // Reserved
		axon.CloseTLSHandshake,     // Reserved
		axon.CloseCode(999),        // Out of range
		axon.CloseCode(2999),       // Out of range
		axon.CloseCode(5000),       // Out of range
	}

	for _, code := range valid {
		if !code.IsValid() {
			t.Errorf("CloseCode(%d).IsValid() = false, want true", code)
		}
	}

	for _, code := range invalid {
		if code.IsValid() {
			t.Errorf("CloseCode(%d).IsValid() = true, want false", code)
		}
	}
}

func TestCloseCodeIsRecoverable(t *testing.T) {
	recoverable := []axon.CloseCode{
		axon.CloseGoingAway,
		axon.CloseAbnormalClosure,
		axon.CloseInternalError,
		axon.CloseServiceRestart,
		axon.CloseTryAgainLater,
		axon.CloseBadGateway,
		axon.CloseCode(4000), // Application-defined, assumed recoverable
	}

	notRecoverable := []axon.CloseCode{
		axon.CloseNormalClosure,
		axon.CloseProtocolError,
		axon.CloseUnsupportedData,
		axon.CloseInvalidPayloadData,
		axon.ClosePolicyViolation,
		axon.CloseMessageTooBig,
		axon.CloseMandatoryExtension,
	}

	for _, code := range recoverable {
		if !code.IsRecoverable() {
			t.Errorf("CloseCode(%d).IsRecoverable() = false, want true", code)
		}
	}

	for _, code := range notRecoverable {
		if code.IsRecoverable() {
			t.Errorf("CloseCode(%d).IsRecoverable() = true, want false", code)
		}
	}
}

func TestCloseError(t *testing.T) {
	t.Run("with reason", func(t *testing.T) {
		err := axon.NewCloseError(4014, "Disallowed intents")
		if err.Code != 4014 {
			t.Errorf("Code = %d, want 4014", err.Code)
		}
		if err.Reason != "Disallowed intents" {
			t.Errorf("Reason = %q, want %q", err.Reason, "Disallowed intents")
		}
		expected := "axon: connection closed (code: 4014 - Application-Defined, reason: Disallowed intents)"
		if got := err.Error(); got != expected {
			t.Errorf("Error() = %q, want %q", got, expected)
		}
	})

	t.Run("without reason", func(t *testing.T) {
		err := axon.NewCloseError(1000, "")
		expected := "axon: connection closed (code: 1000 - Normal Closure)"
		if got := err.Error(); got != expected {
			t.Errorf("Error() = %q, want %q", got, expected)
		}
	})
}

func TestAsCloseError(t *testing.T) {
	t.Run("is CloseError", func(t *testing.T) {
		err := axon.NewCloseError(1000, "test")
		closeErr := axon.AsCloseError(err)
		if closeErr == nil {
			t.Error("AsCloseError() returned nil for CloseError")
		}
		if closeErr.Code != 1000 {
			t.Errorf("Code = %d, want 1000", closeErr.Code)
		}
	})

	t.Run("is not CloseError", func(t *testing.T) {
		closeErr := axon.AsCloseError(axon.ErrConnectionClosed)
		if closeErr != nil {
			t.Error("AsCloseError() returned non-nil for non-CloseError")
		}
	})

	t.Run("nil error", func(t *testing.T) {
		closeErr := axon.AsCloseError(nil)
		if closeErr != nil {
			t.Error("AsCloseError(nil) returned non-nil")
		}
	})
}

func TestIsCloseError(t *testing.T) {
	if !axon.IsCloseError(axon.NewCloseError(1000, "")) {
		t.Error("IsCloseError() = false for CloseError")
	}

	if axon.IsCloseError(axon.ErrConnectionClosed) {
		t.Error("IsCloseError() = true for non-CloseError")
	}
}

