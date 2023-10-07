package dctx

import (
	"context"
)

/*
New
  - If resulting context has no go id, WithNewGoID() is used;
*/
func New(options ...Option) context.Context {
	var ctx = context.Background()
	ctx = OptionWith(ctx, options...)

	if GoID(ctx) == "" {
		ctx = WithNewGoID(ctx)
	}

	return ctx
}
