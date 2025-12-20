package axon

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// queuedMessage represents a message waiting to be sent
type queuedMessage[T any] struct {
	msg     T
	ctx     context.Context
	errCh   chan error
	timeout time.Time
}

// MessageQueue manages queuing of messages during disconnection
type MessageQueue[T any] struct {
	mu       sync.Mutex
	queue    []queuedMessage[T]
	maxSize  int
	timeout  time.Duration
	dropped  atomic.Int64
	enqueued atomic.Int64
	sent     atomic.Int64
	closed   atomic.Bool
}

// newMessageQueue creates a new message queue
func newMessageQueue[T any](maxSize int, timeout time.Duration) *MessageQueue[T] {
	if maxSize <= 0 {
		maxSize = 100
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &MessageQueue[T]{
		queue:   make([]queuedMessage[T], 0, maxSize),
		maxSize: maxSize,
		timeout: timeout,
	}
}

// Enqueue adds a message to the queue
// Returns an error channel that will receive the send result
func (mq *MessageQueue[T]) Enqueue(ctx context.Context, msg T) (chan error, error) {
	if mq.closed.Load() {
		return nil, ErrQueueClosed
	}

	mq.mu.Lock()
	defer mq.mu.Unlock()

	// Check if queue is full
	if len(mq.queue) >= mq.maxSize {
		mq.dropped.Add(1)
		return nil, ErrQueueFull
	}

	errCh := make(chan error, 1)
	qm := queuedMessage[T]{
		msg:     msg,
		ctx:     ctx,
		errCh:   errCh,
		timeout: time.Now().Add(mq.timeout),
	}

	mq.queue = append(mq.queue, qm)
	mq.enqueued.Add(1)

	return errCh, nil
}

// Flush sends all queued messages using the provided send function
func (mq *MessageQueue[T]) Flush(sendFn func(context.Context, T) error) {
	mq.mu.Lock()
	queue := mq.queue
	mq.queue = make([]queuedMessage[T], 0, mq.maxSize)
	mq.mu.Unlock()

	now := time.Now()
	for _, qm := range queue {
		// Check if message has expired
		if now.After(qm.timeout) {
			qm.errCh <- ErrQueueTimeout
			close(qm.errCh)
			mq.dropped.Add(1)
			continue
		}

		// Check if context was cancelled
		if qm.ctx != nil && qm.ctx.Err() != nil {
			qm.errCh <- qm.ctx.Err()
			close(qm.errCh)
			mq.dropped.Add(1)
			continue
		}

		// Send the message
		err := sendFn(qm.ctx, qm.msg)
		qm.errCh <- err
		close(qm.errCh)

		if err == nil {
			mq.sent.Add(1)
		} else {
			mq.dropped.Add(1)
		}
	}
}

// Clear discards all queued messages
func (mq *MessageQueue[T]) Clear() {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for _, qm := range mq.queue {
		qm.errCh <- ErrQueueCleared
		close(qm.errCh)
	}
	mq.dropped.Add(int64(len(mq.queue)))
	mq.queue = mq.queue[:0]
}

// Close closes the queue and discards all pending messages
func (mq *MessageQueue[T]) Close() {
	if mq.closed.Swap(true) {
		return // Already closed
	}
	mq.Clear()
}

// Size returns the current number of queued messages
func (mq *MessageQueue[T]) Size() int {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.queue)
}

// Stats returns queue statistics
func (mq *MessageQueue[T]) Stats() MessageQueueStats {
	return MessageQueueStats{
		CurrentSize: mq.Size(),
		MaxSize:     mq.maxSize,
		Dropped:     mq.dropped.Load(),
		Enqueued:    mq.enqueued.Load(),
		Sent:        mq.sent.Load(),
	}
}

// MessageQueueStats contains queue statistics
type MessageQueueStats struct {
	CurrentSize int
	MaxSize     int
	Dropped     int64
	Enqueued    int64
	Sent        int64
}
