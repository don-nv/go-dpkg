package dctx

import "context"

type Option func(ctx context.Context) context.Context

// With - applies `options` if any to the `ctx` and returns the latest child, else returns unmodified `ctx`.
func With(ctx context.Context, options ...Option) context.Context {
	for _, option := range options {
		ctx = option(ctx)
	}

	return ctx
}

// WithNewGoID - is the same as AddNewGoID, but an Option.
func WithNewGoID() Option {
	return WithGoID(newID())
}

// WithGoID - is the same as AddGoID, but an Option.
func WithGoID(id string) Option {
	return func(ctx context.Context) context.Context {
		return AddGoID(ctx, id)
	}
}

// WithNewXRequestID - is the same as AddNewXRequestID, but an Option.
func WithNewXRequestID() Option {
	return WithXRequestID(newID())
}

// WithXRequestID - is the same as AddXRequestID, but an Option.
func WithXRequestID(id string) Option {
	return func(ctx context.Context) context.Context {
		return AddXRequestID(ctx, id)
	}
}
