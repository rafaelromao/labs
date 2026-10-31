// Package scheduler implements bounded indexed deadlines and an optional owned driver.
package scheduler

import (
	"context"
	"errors"
	"github.com/rafaelromao/labs/solutions/go/algorithms"
	"math"
	"sync"
	"time"
)

type task struct {
	due int64
	run func()
}
type Scheduler struct {
	mu       sync.Mutex
	capacity int
	tasks    map[string]task
	queue    *algorithms.IndexedQueue
	wake     chan struct{}
}

func New(capacity int) *Scheduler {
	if capacity <= 0 {
		panic("positive capacity required")
	}
	return &Scheduler{capacity: capacity, tasks: make(map[string]task), queue: algorithms.NewQueue(), wake: make(chan struct{}, 1)}
}
func (s *Scheduler) notify() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}
func (s *Scheduler) Schedule(id string, due int64, run func()) error {
	if run == nil {
		return errors.New("scheduler: nil callback")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok && len(s.tasks) == s.capacity {
		return errors.New("scheduler: capacity")
	}
	s.tasks[id] = task{due, run}
	s.queue.Put(id, due)
	s.notify()
	return nil
}
func (s *Scheduler) Cancel(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.queue.Remove(id) {
		return false
	}
	delete(s.tasks, id)
	s.notify()
	return true
}
func (s *Scheduler) Next() (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.queue.Peek()
	return p.Value, ok
}

// RunDue claims a finite ordered snapshot before callbacks. A callback panic
// propagates; later callbacks in that claimed snapshot are not retried.
func (s *Scheduler) RunDue(now int64) int {
	s.mu.Lock()
	work := []task{}
	for {
		p, ok := s.queue.Peek()
		if !ok || p.Value > now {
			break
		}
		s.queue.Pop()
		work = append(work, s.tasks[p.Key])
		delete(s.tasks, p.Key)
	}
	s.mu.Unlock()
	for _, t := range work {
		t.run()
	}
	return len(work)
}

// Run owns callback execution until ctx cancellation. now must return monotonic
// nanoseconds at normal clock rate. Use one driver and join its return on shutdown.
// Cancellation cannot recall a claimed task or interrupt an uncooperative callback.
func (s *Scheduler) Run(ctx context.Context, now func() int64) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		s.RunDue(now())
		due, ok := s.Next()
		if !ok {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-s.wake:
			}
			continue
		}
		current := now()
		if due <= current {
			continue
		}
		delay := int64(math.MaxInt64)
		if current >= 0 || due <= math.MaxInt64+current {
			delay = due - current
		}
		timer := time.NewTimer(time.Duration(delay))
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-s.wake:
			timer.Stop()
		case <-timer.C:
		}
	}
}
