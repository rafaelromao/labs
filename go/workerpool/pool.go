// Package workerpool provides bounded admission and drain-on-close workers.
package workerpool

import (
	"context"
	"errors"
)

var ErrClosed = errors.New("workerpool: closed")

var ErrFull = errors.New("workerpool: full")

var ErrAborted = errors.New("workerpool: aborted before start")

// Pool bounds both workers and waiting jobs. Never copy a Pool after use.
type Pool struct{}

func New(workers, capacity int) (*Pool, error) {
	panic("TODO: implement")
}

// Submit rejects immediately when the waiting queue is full or admission closed.
// The buffered result channel never requires the caller to receive for progress.
func (p *Pool) Submit(run func() error) (<-chan error, error) {
	panic("TODO: implement")
}

// Close stops admission and returns immediately; accepted tasks drain.
func (p *Pool) Close() {
	panic("TODO: implement")
}

// Done closes after every accepted task finishes. Uncooperative tasks can delay it forever.
func (p *Pool) Done() <-chan struct{} {
	panic("TODO: implement")
}

// SubmitContext gives cooperative tasks the pool abort signal. The start boundary
// is the worker's protected abort check, not the first instruction of user code.
func (p *Pool) SubmitContext(run func(context.Context) error) (<-chan error, error) {
	panic("TODO: implement")
}

// Abort closes admission, skips work not yet started, and requests cancellation
// of running context-aware tasks. Join Done; arbitrary callbacks cannot be killed.
func (p *Pool) Abort() {
	panic("TODO: implement")
}
