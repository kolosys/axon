package mock

import "github.com/kolosys/axon"

// ErrorInjectionHelpers provides utilities for error injection in tests.
// Error injection methods are available directly on MockConn[T] and MockClient[T]:
//
//	mc := mock.NewMockConn[string](nil)
//	mc.InjectReadError(axon.ErrConnectionClosed)
//
//	_, err := mc.Read(ctx)
//	// err will be axon.ErrConnectionClosed

// CommonErrors provides frequently-used errors for injection
type CommonErrors struct {
	ConnectionClosed error
	ContextCanceled  error
	QueueFull        error
	InvalidState     error
}

// DefaultErrors returns a set of common errors for injection
func DefaultErrors() *CommonErrors {
	return &CommonErrors{
		ConnectionClosed: axon.ErrConnectionClosed,
		ContextCanceled:  axon.ErrContextCanceled,
		QueueFull:        axon.ErrQueueFull,
		InvalidState:     axon.ErrInvalidState,
	}
}

// ErrorSequence allows injecting different errors for successive operations
type ErrorSequence struct {
	errors []error
	index  int
}

// NewErrorSequence creates a new error sequence
func NewErrorSequence(errs ...error) *ErrorSequence {
	return &ErrorSequence{
		errors: errs,
	}
}

// Next returns the next error in the sequence, or nil if exhausted
func (es *ErrorSequence) Next() error {
	if es.index >= len(es.errors) {
		return nil
	}
	err := es.errors[es.index]
	es.index++
	return err
}

// Reset resets the sequence to the beginning
func (es *ErrorSequence) Reset() {
	es.index = 0
}

// ErrorMatcher is a function that matches errors for assertions
type ErrorMatcher func(error) bool

// IsConnectionClosedError returns true if the error is a connection closed error
func IsConnectionClosedError(err error) bool {
	return err == axon.ErrConnectionClosed
}

// IsContextCancelledError returns true if the error is a context canceled error
func IsContextCancelledError(err error) bool {
	return err == axon.ErrContextCanceled
}

// IsQueueFullError returns true if the error is a queue full error
func IsQueueFullError(err error) bool {
	return err == axon.ErrQueueFull
}
