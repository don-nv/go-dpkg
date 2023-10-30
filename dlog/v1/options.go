package dlog

import (
	"context"
)

type Option func(l *Logger)

// WithReadScopeDisabled - disables Log.Scope() and Data.Scope() methods.
func WithReadScopeDisabled() Option {
	return func(l *Logger) {
		l.readScope = func(_ context.Context, l Logger) Data { return l.With() }
	}
}

// WithReadScope - adds custom ReadScopeFn replacing default one for Log.Scope() and Data.Scope() methods.
func WithReadScope(f ReadScopeFn) Option {
	return func(l *Logger) {
		l.readScope = f
	}
}

// WithLevel - sets Logger Level to `lvl`.
func WithLevel(lvl Level) Option {
	return func(l *Logger) {
		l.levels = lvl
	}
}
