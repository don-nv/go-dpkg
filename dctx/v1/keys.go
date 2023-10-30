package dctx

import (
	"context"
	"github.com/google/uuid"
)

func newID() string {
	return uuid.NewString()
}

type keyGoID struct{}

// WithNewGoID - creates `ctx` child, adds respective value and returns it.
func WithNewGoID(ctx context.Context) context.Context {
	return WithGoID(ctx, newID())
}

// WithGoID - creates `ctx` child, adds respective value and returns it.
func WithGoID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyGoID{}, id)
}

func GoID(ctx context.Context) string {
	var v, _ = ctx.Value(keyGoID{}).(string)

	return v
}

type keyXRequestID struct{}

// WithNewXRequestID - creates `ctx` child, adds respective value and returns it.
func WithNewXRequestID(ctx context.Context) context.Context {
	return WithXRequestID(ctx, newID())
}

// WithXRequestID - creates `ctx` child, adds respective value and returns it.
func WithXRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyXRequestID{}, id)
}

func XRequestID(ctx context.Context) string {
	var v, _ = ctx.Value(keyXRequestID{}).(string)

	return v
}
