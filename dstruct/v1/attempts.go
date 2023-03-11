package dstruct

import (
	"context"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"github.com/don-nv/go-dpkg/dtime/v1"
	"time"
)

// AttemptsV1Sync - is the same as AttemptsV1, but it is safe to be used concurrently.
type AttemptsV1Sync struct {
	mu       dsync.RWMutex
	attempts AttemptsV1
}

func NewAttemptsV1Sync(attempts AttemptsV1) AttemptsV1Sync {
	return AttemptsV1Sync{
		attempts: attempts,
	}
}

// Reset - resets attempts.
func (a *AttemptsV1Sync) Reset() {
	a.mu.LockF(func() {
		a.attempts.Reset()
	})
}

// Next - reports whether there is a next attempt. First invocation should be done before actual attempt action.
func (a *AttemptsV1Sync) Next() bool {
	defer a.mu.Lock().Unlock()

	return a.attempts.Next()
}

// AttemptsN - returns attempts number total.
func (a *AttemptsV1Sync) AttemptsN() int {
	defer a.mu.RLock().RUnlock()

	return a.attempts.AttemptsN()
}

// AttemptN - returns current attempt number (not index) among all attempts.
func (a *AttemptsV1Sync) AttemptN() int {
	defer a.mu.RLock().RUnlock()

	return a.attempts.AttemptN()
}

// Exceeded - indicates if all attempts were exceeded without moving to next attempt comparing to Next().
func (a *AttemptsV1Sync) Exceeded() bool {
	defer a.mu.RLock().RUnlock()

	return a.attempts.Exceeded()
}

// AwaitDelay - awaits delay between current and next attempt. If 'ctx' ends before delay, context error is returned.
func (a *AttemptsV1Sync) AwaitDelay(ctx context.Context) error {
	return dtime.AwaitDelay(ctx, a.Delay())
}

// Delay - returns delay between current and next attempts. Any negative delay found is transformed to positive.
func (a *AttemptsV1Sync) Delay() time.Duration {
	defer a.mu.RLock().RUnlock()

	return a.attempts.Delay()
}

type AttemptsV1 struct {
	firstAttempted bool
	delayI         int
	Delays         []time.Duration `json:"delays" yaml:"delays"`
}

func (a *AttemptsV1) hasDelays() bool {
	return len(a.Delays) > 0
}

func (a *AttemptsV1) lastDelayI() int {
	return len(a.Delays) - 1
}

// Reset - resets attempts.
func (a *AttemptsV1) Reset() {
	a.firstAttempted = false
	a.delayI = 0
}

// Next - reports whether there is a next attempt. First invocation should be done before actual attempt action.
func (a *AttemptsV1) Next() bool {
	if !a.firstAttempted {
		a.firstAttempted = true

		return a.hasDelays()
	}

	i := a.delayI + 1
	if i > a.lastDelayI() {
		return false
	}

	a.delayI = i

	return true
}

// AttemptsN - returns attempts number total.
func (a *AttemptsV1) AttemptsN() int {
	return len(a.Delays)
}

// AttemptN - returns current attempt number (not index) among all attempts.
func (a *AttemptsV1) AttemptN() int {
	return a.delayI + 1
}

// Exceeded - indicates if all attempts were exceeded without moving to next attempt comparing to Next().
func (a *AttemptsV1) Exceeded() bool {
	return a.AttemptN() >= a.AttemptsN()
}

// AwaitDelay - awaits delay between current and next attempt. If 'ctx' ends before delay, context error is returned.
func (a *AttemptsV1) AwaitDelay(ctx context.Context) error {
	return dtime.AwaitDelay(ctx, a.Delay())
}

// Delay - returns delay after current attempt. Negative durations become positive and 0 become 1.
func (a *AttemptsV1) Delay() time.Duration {
	delay := a.Delays[a.delayI]
	if delay < 0 {
		return delay * -1
	}
	if delay < 1 {
		return 1
	}

	return delay
}

type AttemptsSync struct {
	mu       dsync.RWMutex
	attempts Attempts
}

func NewAttemptsSync(attempts Attempts) AttemptsSync {
	return AttemptsSync{
		attempts: attempts,
	}
}

func (a *AttemptsSync) Exceeded() bool {
	defer a.mu.RLock().RUnlock()

	return a.attempts.Exceeded()
}

func (a *AttemptsSync) Next() bool {
	defer a.mu.Lock().Unlock()

	return a.attempts.Next()
}

func (a *AttemptsSync) Reset() {
	a.mu.LockF(func() {
		a.attempts.Reset()
	})
}

func (a *AttemptsSync) AtLeastOnceAttempted() bool {
	defer a.mu.RLock().RUnlock()

	return a.attempts.AtLeastOnceAttempted()
}

func (a *AttemptsSync) AttemptN() int {
	defer a.mu.RLock().RUnlock()

	return a.attempts.AttemptN()
}

func (a *AttemptsSync) AttemptsN() int {
	defer a.mu.RLock().RUnlock()

	return a.attempts.AttemptsN()
}

type Attempts struct {
	MaxN     int `json:"max_n" yaml:"max_n"`
	currentN int
}

func (a *Attempts) AtLeastOnceAttempted() bool {
	return a.currentN > 0
}

func (a *Attempts) Exceeded() bool {
	return a.currentN == a.MaxN
}

func (a *Attempts) Next() bool {
	n := a.currentN + 1
	if n > a.MaxN {
		return false
	}

	a.currentN = n

	return true
}

func (a *Attempts) Reset() {
	a.currentN = 0
}

func (a *Attempts) AttemptN() int {
	return a.currentN
}

func (a *Attempts) AttemptsN() int {
	return a.MaxN
}
