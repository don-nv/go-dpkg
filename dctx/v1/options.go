package dctx

import "context"

type Option func(ctx context.Context) context.Context

func With(ctx context.Context, options ...Option) context.Context {
	for _, option := range options {
		ctx = option(ctx)
	}

	return ctx
}

func WithNewGoID() Option {
	return WithGoID(newID())
}

func WithGoID(id string) Option {
	return func(ctx context.Context) context.Context {
		return AddGoID(ctx, id)
	}
}

func WithNewXRequestID() Option {
	return WithXRequestID(newID())
}

func WithXRequestID(id string) Option {
	return func(ctx context.Context) context.Context {
		return AddXRequestID(ctx, id)
	}
}
