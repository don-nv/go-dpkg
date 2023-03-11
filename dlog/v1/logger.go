package dlog

import (
	"github.com/don-nv/go-dpkg/derr/v1"
	"github.com/rs/zerolog"
	"os"
	"time"
)

/*
Logger - writes logs. This is a generalized implementation having well-known general and custom-minded convenient API.
*/
type Logger struct {
	zero       zerolog.Logger
	startedAt  time.Time
	readScope  ReadScopeFn
	names      [][]byte
	catchEDMsg string
	levels     Level
}

/*
New
  - Level is set to LevelAll by default;
  - LevelDebug message logged at Logger.CatchED() method is CatchEDDefaultMessage by default;
  - ReadScopeDefault is used by default;
*/
func New(options ...OptionLogger) Logger {
	log := Logger{
		levels:     LevelAll,
		readScope:  ReadScopeDefault,
		catchEDMsg: CatchEDDefaultMessage,
	}

	for _, option := range options {
		option(&log)
	}

	log.zero = zerolog.New(os.Stdout)

	return log
}

// E - returns new Log at LevelError.
func (l Logger) E() Log { return l.newLog(LevelError) }

// W - returns new Log at LevelWarn.
func (l Logger) W() Log { return l.newLog(LevelWarn) }

// WithDuration - enables duration writing passed since WithDuration call and actual log writing.
func (l Logger) WithDuration() Logger { l.startedAt = time.Now(); return l }

// I - returns new Log at LevelInfo.
func (l Logger) I() Log { return l.newLog(LevelInfo) }

// D - returns new Log at LevelDebug.
func (l Logger) D() Log { return l.newLog(LevelDebug) }

// newLog - returns new Log. If 'lvl' is disabled, the returned Log is no-op.
func (l Logger) newLog(lvl Level) Log {
	// Preconfigure zerolog. If disabled, does nothing.
	if !l.levels.Enabled(lvl) {
		l.zero = l.zero.Level(zerolog.Disabled)
	}

	return newLog(l, lvl)
}

// With - returns Logger Data to be populated. Call Data.Build() to return a Logger with new data added.
func (l Logger) With() Data {
	return newData(l)
}

/*
CatchE - catches an 'err' and if '*err' != nil, writes Log at E() method. If 'notErrs' are passed, then '*err' gets
compared against them and if found, nothing will be logged as well as if '*err' == nil. This should be used with 'defer`
directive. It reveals operation intermediaries if an error occurred. However, this method might be tricky to use in
cases when it comes to error shadowing.

Example:

	func F(log Logger) {
		var err error
		defer log.CatchE(&err)
		...
		log = log.With().Any("k", "v").Build()
		err = errors.New("an error")
		...
		// [E] an error {"k":"v"}
	}
*/
func (l *Logger) CatchE(err *error, notErrs ...error) {
	// CatchE, CatchED methods must point to *Logger. Otherwise, data added after 'defer log.CatchED()' get lost.
	l.catchED(false, err, notErrs...)
}

/*
CatchED - is the same as CatchE() method, but writes Log at D() method with predefined message if '*err' == nil or
there's no match with 'notErrs'.

Example:

	func F(log Logger) {
		var err error
		defer log.CatchED(&err)
		...
		log = log.With().Any("k", "v").Build()
		...
		// [D] OK {"k":"v"}
	}
*/
func (l *Logger) CatchED(err *error, notErrs ...error) {
	// CatchE, CatchED methods must point to *Logger. Otherwise, data added after 'defer log.CatchED()' get lost.
	l.catchED(true, err, notErrs...)
}

func (l Logger) catchED(isDebugCatch bool, err *error, notErrs ...error) {
	if *err == nil {
		if isDebugCatch {
			l.D().Write(l.catchEDMsg)
		}

		return
	}

	ok := derr.InP(err, notErrs...)
	if !ok {
		l.E().Write((*err).Error())

		return
	}

	if isDebugCatch {
		l.D().Write(l.catchEDMsg)
	}
}

func (l Logger) constructName() string {
	if len(l.names) < 1 {
		return ""
	}

	// TODO? Cap may be replaced with bytes sync pools.
	var name = make([]byte, 0, len(l.names)*nameExpectedMaxBytes)
	for i, n := range l.names {
		name = append(name, n...)

		if i != len(l.names)-1 {
			name = append(name, '.')
		}
	}

	return string(name)
}
