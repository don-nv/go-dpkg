package dctx_test

import (
	"context"
	"dpkg/dctx"
	"testing"
)

func TestTODO(t *testing.T) {
	ctx := dctx.WithValues(context.Background(), dctx.WithNewRequestID, dctx.WithNewGoID)

	t.Log(dctx.RequestID(ctx))
	t.Log(dctx.GoID(ctx))

	ctx2 := dctx.Detach(ctx,
		func(ctx context.Context) (k, v interface{}) {
			return dctx.KeyGoID{}, dctx.GoID(ctx)
		},
		func(ctx context.Context) (k, v interface{}) {
			return dctx.KeyXRequestID{}, dctx.RequestID(ctx)
		},
	)

	t.Log(dctx.RequestID(ctx2))
	t.Log(dctx.GoID(ctx2))
}
