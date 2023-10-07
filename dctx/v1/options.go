package dctx

import "context"

type Option func(ctx context.Context) context.Context

func OptionWith(ctx context.Context, options ...Option) context.Context {
	for _, option := range options {
		ctx = option(ctx)
	}

	return ctx
}

func OptionWithNewGoID() Option {
	return OptionWithGoID(NewID())
}

func OptionWithGoID(id string) Option {
	return func(ctx context.Context) context.Context {
		return WithGoID(ctx, id)
	}
}

func OptionWithNewXRequestID() Option {
	return OptionWithXRequestID(NewID())
}

func OptionWithXRequestID(id string) Option {
	return func(ctx context.Context) context.Context {
		return WithXRequestID(ctx, id)
	}
}
