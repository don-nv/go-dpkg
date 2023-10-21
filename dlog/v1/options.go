package dlog

import (
	"context"
)

type Option func(l *Logger)

// WithReadScopeDisabled - disables Log.Scope() and Builder.Scope() methods.
func WithReadScopeDisabled() Option {
	return func(l *Logger) {
		l.readCtx = func(_ context.Context, l Logger) Logger { return l }
	}
}

// WithReadScope - adds custom ReadScopeFn replacing default one for Log.Scope() and Builder.Scope() methods.
func WithReadScope(f ReadScopeFn) Option {
	return func(l *Logger) {
		l.readCtx = f
	}
}

// WithLevel - sets Logger Level to `lvl`.
func WithLevel(lvl Level) Option {
	return func(l *Logger) {
		l.levels = lvl
	}
}
