package dctx

import (
	"context"
	"github.com/don-nv/go-dpkg/dchan/v1"
	"os"
	"os/signal"
	"time"
)

// NewTTL - is the same as New, but returns context with a limited time-to-live.
func NewTTL(options ...TTLOption) (context.Context, context.CancelFunc) {
	return WithTTLOptions(New(), options...)
}

/*
TTLWithOSCancel - creates [ctx] child with a limited time-to-live dependent on os.Interrupt and os.Kill signals.
  - 'signals' - are optional additional signals specific for the OS;
*/
func TTLWithOSCancel(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	var s = append(make([]os.Signal, 0, len(signals)+2),
		os.Interrupt,
		os.Kill,
	)
	return signal.NotifyContext(
		ctx,
		append(s, signals...)...,
	)
}

// WithTTLCancel - is a regular context.WithCancel wraparound.
func WithTTLCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

// WithTTLTimeout - is a regular context.WithTimeout wraparound.
func WithTTLTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}

// WithTTLC - creates [ctx] child and returns it. Child context is canceled when 'c' sends a signal or [ctx] is done.
func WithTTLC(ctx context.Context, c <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := WithTTLCancel(ctx)

	go func() {
		dchan.Receive(ctx, c)

		cancel()
	}()

	return ctx, cancel
}

/*
Timeout - returns [ctx] Timeout, a duration relative to current time before, after [ctx] expires.
  - < 0: expired
  - = 0: no Timeout
  - > 0: not expired
*/
func Timeout(ctx context.Context) time.Duration {
	t, ok := ctx.Deadline()
	if !ok {
		return 0
	}

	d := time.Until(t)
	if d == 0 {
		return -1
	}

	return d
}
