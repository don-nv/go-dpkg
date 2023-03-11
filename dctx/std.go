package dctx

import "context"

func Detach(ctx context.Context, readCtxs ...func(ctx context.Context) (k, v interface{})) context.Context {
	var newCtx = context.Background()

	for _, readCtx := range readCtxs {
		k, v := readCtx(ctx)

		newCtx = WithValues(
			ctx,
			func(ctx context.Context) context.Context {
				return context.WithValue(ctx, k, v)
			},
		)
	}

	return newCtx
}
