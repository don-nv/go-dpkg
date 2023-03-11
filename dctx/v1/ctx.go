package dctx

import (
	"context"
)

/*
New
  - If resulting context has no go id, WithNewGoID gets applied;
*/
func New(options ...Option) context.Context {
	var ctx = WithOptions(context.Background(), options...)

	if GoID(ctx) == "" {
		ctx = WithNewGoID(ctx)
	}

	return ctx
}

// WithoutCancel - is a regular context.WithoutCancel wraparound.
func WithoutCancel(ctx context.Context) context.Context {
	return context.WithoutCancel(ctx)
}
