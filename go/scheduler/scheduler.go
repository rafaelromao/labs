// Package scheduler implements bounded indexed deadlines and an optional owned driver.
package scheduler

import (
	"context"
)

type Scheduler struct{}

func New(capacity int) *Scheduler {
	panic("TODO: implement")
}

func (s *Scheduler) Schedule(id string, due int64, run func()) error {
	panic("TODO: implement")
}

func (s *Scheduler) Cancel(id string) bool {
	panic("TODO: implement")
}

func (s *Scheduler) Next() (int64, bool) {
	panic("TODO: implement")
}

// RunDue claims a finite ordered snapshot before callbacks. A callback panic
// propagates; later callbacks in that claimed snapshot are not retried.
func (s *Scheduler) RunDue(now int64) int {
	panic("TODO: implement")
}

// Run owns callback execution until ctx cancellation. now must return monotonic
// nanoseconds at normal clock rate. Use one driver and join its return on shutdown.
// Cancellation cannot recall a claimed task or interrupt an uncooperative callback.
func (s *Scheduler) Run(ctx context.Context, now func() int64) error {
	panic("TODO: implement")
}
