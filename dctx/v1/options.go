package dctx

import "context"

type Option func(ctx context.Context) context.Context

// WithOptions - applies `options` if any to the `ctx` and returns the latest child, else returns unmodified `ctx`.
func WithOptions(ctx context.Context, options ...Option) context.Context {
	for _, option := range options {
		ctx = option(ctx)
	}

	return ctx
}

// OptionNewGoID - is the same as WithNewGoID, but an Option.
func OptionNewGoID() Option {
	return OptionGoID(newID())
}

// OptionGoID - is the same as WithGoID, but an Option.
func OptionGoID(id string) Option {
	return func(ctx context.Context) context.Context {
		return WithGoID(ctx, id)
	}
}

// OptionNewXRequestID - is the same as WithNewXRequestID, but an Option.
func OptionNewXRequestID() Option {
	return OptionXRequestID(newID())
}

// OptionXRequestID - is the same as WithXRequestID, but an Option.
func OptionXRequestID(id string) Option {
	return func(ctx context.Context) context.Context {
		return WithXRequestID(ctx, id)
	}
}
