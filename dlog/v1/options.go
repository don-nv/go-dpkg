package dlog

import (
	"context"
)

type OptionLogger func(l *Logger)

// OptionLoggerWithReadScopeDisabled - disables Log.Scope() and Data.Scope() methods.
func OptionLoggerWithReadScopeDisabled() OptionLogger {
	return func(l *Logger) {
		l.readScope = func(_ context.Context, data Data) Data { return data }
	}
}

// OptionLoggerWithReadScope - adds custom ReadScopeFn replacing default one for Log.Scope() and Data.Scope() methods.
func OptionLoggerWithReadScope(f ReadScopeFn) OptionLogger {
	return func(l *Logger) {
		l.readScope = f
	}
}

// OptionLoggerWithLevel - sets Logger Level to 'lvl'.
func OptionLoggerWithLevel(lvl Level) OptionLogger {
	return func(l *Logger) {
		l.levels = lvl
	}
}
