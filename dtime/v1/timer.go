package dtime

import (
	"context"
	"github.com/don-nv/go-dpkg/dchan/v1"
	"time"
)

type Timer struct {
	timer *time.Timer
}

func NewTimer(d time.Duration) Timer {
	var timer = time.NewTimer(d)

	return Timer{
		timer: timer,
	}
}

func (t Timer) C() <-chan time.Time {
	return t.timer.C
}

// Stop - is the same as time.Timer.Stop(), but drains timers' channel after.
func (t Timer) Stop() bool {
	var ok = t.timer.Stop()

	if !ok {
		dchan.Drain(context.Background(), t.timer.C)
	}

	return ok
}

// Reset - is the same as time.Timer.Reset(), but drains timers' channel after.
func (t Timer) Reset(d time.Duration) bool {
	var ok = t.timer.Reset(d)

	if !ok {
		dchan.Drain(context.Background(), t.timer.C)
	}

	return ok
}
