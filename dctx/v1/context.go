package dctx

import (
	"context"
)

/*
New
  - If resulting context has no go id, OptionNewGoID() gets applied;
*/
func New(options ...Option) context.Context {
	var ctx = WithOptions(context.Background(), options...)

	if GoID(ctx) == "" {
		ctx = OptionNewGoID()(ctx)
	}

	return ctx
}
