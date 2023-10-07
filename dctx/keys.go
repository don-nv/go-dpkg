package dctx

import (
	"context"
	"github.com/google/uuid"
)

func NewID() string {
	return uuid.NewString()
}

type keyGoID struct{}

func WithNewGoID(ctx context.Context) context.Context {
	return WithGoID(ctx, NewID())
}

func WithGoID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyGoID{}, id)
}

func GoID(ctx context.Context) string {
	return ctx.Value(keyGoID{}).(string)
}

type keyXRequestID struct{}

func WithNewXRequestID(ctx context.Context) context.Context {
	return WithXRequestID(ctx, NewID())
}

func WithXRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyXRequestID{}, id)
}

func XRequestID(ctx context.Context) string {
	return ctx.Value(keyXRequestID{}).(string)
}
