package axon

import "fmt"

// CloseCode represents a WebSocket close status code (RFC 6455 Section 7.4).
type CloseCode int

// Standard WebSocket close codes (RFC 6455 Section 7.4.1)
const (
	CloseNormalClosure      CloseCode = 1000 // Normal closure; the connection successfully completed
	CloseGoingAway          CloseCode = 1001 // Endpoint going away (server shutdown, browser navigated away)
	CloseProtocolError      CloseCode = 1002 // Protocol error encountered
	CloseUnsupportedData    CloseCode = 1003 // Received data type cannot be accepted
	CloseNoStatusReceived   CloseCode = 1005 // No status code was present (reserved, must not be sent)
	CloseAbnormalClosure    CloseCode = 1006 // Connection closed abnormally (reserved, must not be sent)
	CloseInvalidPayloadData CloseCode = 1007 // Received data was inconsistent with message type
	ClosePolicyViolation    CloseCode = 1008 // Message violated policy (generic code when 1003/1009 don't apply)
	CloseMessageTooBig      CloseCode = 1009 // Message too large to process
	CloseMandatoryExtension CloseCode = 1010 // Server didn't negotiate required extension
	CloseInternalError      CloseCode = 1011 // Server encountered an unexpected condition
	CloseServiceRestart     CloseCode = 1012 // Server is restarting
	CloseTryAgainLater      CloseCode = 1013 // Server is overloaded, try again later
	CloseBadGateway         CloseCode = 1014 // Server acting as gateway received invalid response
	CloseTLSHandshake       CloseCode = 1015 // TLS handshake failed (reserved, must not be sent)
)

// String returns a human-readable name for the close code.
func (c CloseCode) String() string {
	switch c {
	case CloseNormalClosure:
		return "Normal Closure"
	case CloseGoingAway:
		return "Going Away"
	case CloseProtocolError:
		return "Protocol Error"
	case CloseUnsupportedData:
		return "Unsupported Data"
	case CloseNoStatusReceived:
		return "No Status Received"
	case CloseAbnormalClosure:
		return "Abnormal Closure"
	case CloseInvalidPayloadData:
		return "Invalid Payload Data"
	case ClosePolicyViolation:
		return "Policy Violation"
	case CloseMessageTooBig:
		return "Message Too Big"
	case CloseMandatoryExtension:
		return "Mandatory Extension"
	case CloseInternalError:
		return "Internal Error"
	case CloseServiceRestart:
		return "Service Restart"
	case CloseTryAgainLater:
		return "Try Again Later"
	case CloseBadGateway:
		return "Bad Gateway"
	case CloseTLSHandshake:
		return "TLS Handshake Failed"
	default:
		if c >= 3000 && c < 4000 {
			return "Library-Defined"
		}
		if c >= 4000 && c < 5000 {
			return "Application-Defined"
		}
		return "Unknown"
	}
}

// IsReserved returns true if this is a reserved close code that must not be sent in close frames.
func (c CloseCode) IsReserved() bool {
	switch c {
	case CloseNoStatusReceived, CloseAbnormalClosure, CloseTLSHandshake:
		return true
	default:
		return false
	}
}

// IsValid returns true if this is a valid close code that can be sent in close frames.
func (c CloseCode) IsValid() bool {
	if c.IsReserved() {
		return false
	}
	// Valid ranges: 1000-1015 (standard), 3000-3999 (library), 4000-4999 (application)
	return (c >= 1000 && c <= 1015) || (c >= 3000 && c < 5000)
}

// IsRecoverable returns true if reconnection should typically be attempted.
// This is a hint - applications may have their own logic for specific codes.
func (c CloseCode) IsRecoverable() bool {
	switch c {
	case CloseNormalClosure:
		return false // Intentional close, don't reconnect
	case CloseGoingAway:
		return true // Server shutting down, try again
	case CloseProtocolError, CloseUnsupportedData, CloseInvalidPayloadData:
		return false // Client error, reconnecting won't help
	case ClosePolicyViolation:
		return false // Policy violation, reconnecting won't help
	case CloseMessageTooBig:
		return false // Message too big, client needs to fix
	case CloseMandatoryExtension:
		return false // Missing extension, client needs to fix
	case CloseInternalError:
		return true // Server error, may be transient
	case CloseServiceRestart, CloseTryAgainLater:
		return true // Temporary, should retry
	case CloseBadGateway:
		return true // Gateway error, may be transient
	case CloseAbnormalClosure, CloseNoStatusReceived:
		return true // Connection lost, try again
	default:
		// For application-defined codes (4000+), assume recoverable by default
		// Applications should provide their own ShouldReconnect logic
		return c >= 4000
	}
}

// CloseError represents a WebSocket close event with code and reason.
// It implements the error interface for use in error handling.
type CloseError struct {
	Code   CloseCode
	Reason string
}

// Error implements the error interface.
func (e *CloseError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("axon: connection closed (code: %d - %s, reason: %s)", e.Code, e.Code.String(), e.Reason)
	}
	return fmt.Sprintf("axon: connection closed (code: %d - %s)", e.Code, e.Code.String())
}

// IsRecoverable returns true if reconnection should typically be attempted.
func (e *CloseError) IsRecoverable() bool {
	return e.Code.IsRecoverable()
}

// NewCloseError creates a new CloseError with the given code and reason.
func NewCloseError(code int, reason string) *CloseError {
	return &CloseError{
		Code:   CloseCode(code),
		Reason: reason,
	}
}

// AsCloseError attempts to extract a CloseError from an error.
// Returns nil if the error is not a CloseError.
func AsCloseError(err error) *CloseError {
	if closeErr, ok := err.(*CloseError); ok {
		return closeErr
	}
	return nil
}

// IsCloseError returns true if the error is a CloseError.
func IsCloseError(err error) bool {
	return AsCloseError(err) != nil
}

