package dctx

import (
	"context"
	"github.com/google/uuid"
	"log"
	"runtime/debug"
	"strings"
)

func WithValues(ctx context.Context, sets ...func(ctx context.Context) context.Context) context.Context {
	for _, set := range sets {
		ctx = set(ctx)
	}

	return ctx
}

func getString(ctx context.Context, key interface{}) string {
	switch v := ctx.Value(key).(type) {
	case string:
		return v

	case nil:
		return ""

	default:
		const format = "[WARN] getting string value from context, value has unexpected %T type\n\nstack:\n\n%s\n"
		log.Printf(format, v, debug.Stack())

		return ""
	}
}

// KeyGoID - is a context key that refers to assigned goroutine id value.
type KeyGoID struct{}

func (k KeyGoID) String() string {
	return "go_id"
}

// WithNewGoID - returns new context with new goroutine id value set.
func WithNewGoID(ctx context.Context) context.Context {
	return context.WithValue(ctx, KeyGoID{}, newGoID())
}

func WithGoID(ctx context.Context, goID string) context.Context {
	return context.WithValue(ctx, KeyGoID{}, goID)
}

func newGoID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func GoID(ctx context.Context) string {
	return getString(ctx, KeyGoID{})
}

// KeyXRequestID - is a context key that refers to assigned request id value.
type KeyXRequestID struct{}

func (k KeyXRequestID) String() string {
	return "x_request_id"
}

// WithNewRequestID - returns new context with new request id value set.
func WithNewRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, KeyXRequestID{}, NewRequestID())
}

// WithRequestID - returns new context with passed request id value set.
func WithRequestID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, KeyXRequestID{}, value)
}

func NewRequestID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func RequestID(ctx context.Context) string {
	return getString(ctx, KeyXRequestID{})
}
