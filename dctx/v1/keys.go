package dctx

import (
	"context"
	"github.com/google/uuid"
)

func newID() string {
	return uuid.NewString()
}

type keyGoID struct{}

func AddNewGoID(ctx context.Context) context.Context {
	return AddGoID(ctx, newID())
}

func AddGoID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyGoID{}, id)
}

func GoID(ctx context.Context) string {
	var v, _ = ctx.Value(keyGoID{}).(string)

	return v
}

type keyXRequestID struct{}

func AddNewXRequestID(ctx context.Context) context.Context {
	return AddXRequestID(ctx, newID())
}

func AddXRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyXRequestID{}, id)
}

func XRequestID(ctx context.Context) string {
	var v, _ = ctx.Value(keyXRequestID{}).(string)

	return v
}
