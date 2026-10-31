// Package workerpool provides bounded admission and drain-on-close workers.
package workerpool

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrClosed = errors.New("workerpool: closed")
var ErrFull = errors.New("workerpool: full")
var ErrAborted = errors.New("workerpool: aborted before start")

type job struct {
	run    func() error
	result chan error
}

// Pool bounds both workers and waiting jobs. Never copy a Pool after use.
type Pool struct {
	mu      sync.Mutex
	closed  bool
	aborted bool
	ctx     context.Context
	cancel  context.CancelFunc
	jobs    chan job
	done    chan struct{}
}

func New(workers, capacity int) (*Pool, error) {
	if workers <= 0 || capacity <= 0 {
		return nil, errors.New("workerpool: positive workers and capacity required")
	}
	p := &Pool{jobs: make(chan job, capacity), done: make(chan struct{})}
	p.ctx, p.cancel = context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := range p.jobs {
				p.mu.Lock()
				aborted := p.aborted
				p.mu.Unlock()
				if aborted {
					j.result <- ErrAborted
				} else {
					j.result <- execute(j.run)
				}
				close(j.result)
			}
		}()
	}
	go func() { wg.Wait(); p.cancel(); close(p.done) }()
	return p, nil
}

func execute(run func() error) (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("workerpool: task panicked: %v", value)
		}
	}()
	return run()
}

// Submit rejects immediately when the waiting queue is full or admission closed.
// The buffered result channel never requires the caller to receive for progress.
func (p *Pool) Submit(run func() error) (<-chan error, error) {
	if run == nil {
		return nil, errors.New("workerpool: nil task")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil, ErrClosed
	}
	result := make(chan error, 1)
	select {
	case p.jobs <- job{run, result}:
		return result, nil
	default:
		return nil, ErrFull
	}
}

// Close stops admission and returns immediately; accepted tasks drain.
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.closed {
		p.closed = true
		close(p.jobs)
	}
}

// Done closes after every accepted task finishes. Uncooperative tasks can delay it forever.
func (p *Pool) Done() <-chan struct{} { return p.done }

// SubmitContext gives cooperative tasks the pool abort signal. The start boundary
// is the worker's protected abort check, not the first instruction of user code.
func (p *Pool) SubmitContext(run func(context.Context) error) (<-chan error, error) {
	if run == nil {
		return nil, errors.New("nil task")
	}
	return p.Submit(func() error { return run(p.ctx) })
}

// Abort closes admission, skips work not yet started, and requests cancellation
// of running context-aware tasks. Join Done; arbitrary callbacks cannot be killed.
func (p *Pool) Abort() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.aborted = true
	p.cancel()
	if !p.closed {
		p.closed = true
		close(p.jobs)
	}
}
