package dctx

import (
	"context"
)

/*
New
  - If resulting context has no go id, WithNewGoID() gets applied;
*/
func New(options ...Option) context.Context {
	var ctx = With(context.Background(), options...)

	if GoID(ctx) == "" {
		ctx = WithNewGoID()(ctx)
	}

	return ctx
}
