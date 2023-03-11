package dctx

import (
	"context"
	"os"
	"time"
)

type Option func(ctx context.Context) context.Context

// WithOptions - applies 'options' if any to the 'ctx' and returns the latest child, else returns unmodified 'ctx'.
func WithOptions(ctx context.Context, options ...Option) context.Context {
	for _, option := range options {
		ctx = option(ctx)
	}

	return ctx
}

// OptionWithOSCancel - is the same as OptionTTLWithOSCancel, but discards cancel func.
func OptionWithOSCancel(signals ...os.Signal) Option {
	return func(ctx context.Context) context.Context {
		ctx, _ = OptionTTLWithOSCancel(signals...)(ctx)

		return ctx
	}
}

// OptionWithNewGoID - is the same as WithNewGoID, but an Option.
func OptionWithNewGoID() Option {
	return OptionWithGoID(newID())
}

// OptionWithGoID - is the same as WithGoID, but an Option.
func OptionWithGoID(id string) Option {
	return func(ctx context.Context) context.Context {
		return WithGoID(ctx, id)
	}
}

// OptionWithNewXRequestID - is the same as WithNewXRequestID, but an Option.
func OptionWithNewXRequestID() Option {
	return OptionWithXRequestID(newID())
}

// OptionWithXRequestID - is the same as WithXRequestID, but an Option.
func OptionWithXRequestID(id string) Option {
	return func(ctx context.Context) context.Context {
		return WithXRequestID(ctx, id)
	}
}

// TTLOption - options context with time-to-live condition.
type TTLOption func(ctx context.Context) (context.Context, context.CancelFunc)

/*
WithTTLOptions - applies 'options' if any to the 'ctx' and returns the latest child, else returns cancellable 'ctx'
child.
*/
func WithTTLOptions(ctx context.Context, options ...TTLOption) (context.Context, context.CancelFunc) {
	ctx, cancel := WithTTLCancel(ctx)
	for _, option := range options {
		ctx, cancel = option(ctx)
	}

	return ctx, cancel
}

// OptionTTLWithOSCancel - is the same as TTLWithOSCancel, but an TTLOption.
func OptionTTLWithOSCancel(signals ...os.Signal) TTLOption {
	return func(ctx context.Context) (context.Context, context.CancelFunc) {
		return TTLWithOSCancel(ctx, signals...)
	}
}

// OptionTTLWithC - is the same as WithTTLC, but an option.
func OptionTTLWithC(c <-chan struct{}) TTLOption {
	return func(ctx context.Context) (context.Context, context.CancelFunc) {
		return WithTTLC(ctx, c)
	}
}

func OptionTTLWithTimeout(d time.Duration) TTLOption {
	return func(ctx context.Context) (context.Context, context.CancelFunc) {
		return WithTTLTimeout(ctx, d)
	}
}
