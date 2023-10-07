package dlog

import (
	"github.com/don-nv/go-dpkg/derr/v1"
	"github.com/rs/zerolog"
	"os"
)

/*
Logger - writes logs. This is a generalized implementation having well-known general and custom-minded convenient API.
*/
type Logger struct {
	zero       zerolog.Logger
	names      []string
	catchEDMsg string
	lvl        Level
	readCtx    ReadScopeFn
}

/*
New
  - Level is set to LevelAll by default;
  - LevelDebug message logged at Logger.CatchED() method is CatchEDDefaultMessage by default;
  - ReadScopeDefault is used by default;
*/
func New(options ...Option) Logger {
	log := Logger{
		lvl:        LevelAll,
		readCtx:    ReadScopeDefault,
		catchEDMsg: CatchEDDefaultMessage,
	}

	for _, option := range options {
		option(&log)
	}

	zerolog.LevelErrorValue = LevelError.String()
	zerolog.LevelWarnValue = LevelWarn.String()
	zerolog.LevelInfoValue = LevelInfo.String()
	zerolog.LevelDebugValue = LevelDebug.String()
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = TimeDefaultLayout

	log.zero = zerolog.New(os.Stdout)

	return log
}

// E - returns new Log at LevelError.
func (l Logger) E() Log { return l.newLog(LevelError) }

/*
WriteE - writes `msg` at LevelError and returns Log at respective level. This is used to make instant writes to the
level. Also, it is convenient to be used with `defer` directive.

Example:

	func F() {
		...
		log.WriteE("running...").Write("...done")
		// [E] running...
		...
		return
		// [E] ...done
	}
*/
func (l Logger) WriteE(err error) Log { return l.newLog(LevelError).Write(errStr(err)) }

// W - returns new Log at LevelWarn.
func (l Logger) W() Log { return l.newLog(LevelWarn) }

/*
WriteW - writes `msg` at LevelWarn and returns Log at respective level. This is used to make instant writes to the
level. Also, it is convenient to be used with `defer` directive.

Example:

	func F() {
		...
		log.WriteW("running...").Write("...done")
		// [W] running...
		...
		return
		// [W] ...done
	}
*/
func (l Logger) WriteW(msg string) Log { return l.newLog(LevelWarn).Write(msg) }

// I - returns new Log at LevelInfo.
func (l Logger) I() Log { return l.newLog(LevelInfo) }

/*
WriteI - writes `msg` at LevelInfo and returns Log at respective level. This is used to make instant writes to the
level. Also, it is convenient to be used with `defer` directive.

Example:

	func F() {
		...
		log.WriteI("running...").Write("...done")
		// [I] running...
		...
		return
		// [I] ...done
	}
*/
func (l Logger) WriteI(msg string) Log { return l.newLog(LevelInfo).Write(msg) }

// D - returns new Log at LevelDebug.
func (l Logger) D() Log { return l.newLog(LevelDebug) }

/*
WriteD - writes `msg` at LevelDebug and returns Log at respective level. This is used to make instant writes to the
level. Also, it is convenient to be used with `defer` directive.

Example:

	func F() {
		...
		log.WriteD("running...").Write("...done")
		// [D] running...
		...
		return
		// [D] ...done
	}
*/
func (l Logger) WriteD(msg string) Log { return l.newLog(LevelDebug).Write(msg) }

// newLog - returns new Log. If `lvl` is disabled, the returned Log is no-op.
func (l Logger) newLog(lvl Level) Log {
	if !l.lvl.Enabled(lvl) {
		l.zero = l.zero.Level(zerolog.Disabled)
	}

	return newLog(l, lvl)
}

// With - returns Logger Data to be populated. Call Data.Build() to return a Logger with new data added.
func (l Logger) With() Data { return newData(l) }

/*
CatchE - catches an `err` and if `*err` != nil, writes Log at E() method. If `notErrs` are passed, then `*err` gets
compared against them and if found, nothing will be logged as well as if `*err` == nil. This should be used with `defer`
directive. It reveals operation intermediaries if an error occurred. However, this method might be tricky to use in
cases when it comes to error shadowing.

Example:

	func F(log Logger) {
		var err error
		defer log.CatchE(&err)
		...
		log = log.With().Any("k", "v").Build()
		...
		err = errors.New("an error occurred")
		// [E] an error occurred {"k":"v"}
	}
*/
func (l *Logger) CatchE(err *error, notErrs ...error) { l.catchED(false, err, notErrs...) }

/*
CatchED - is the same as CatchE() method, but writes Log at D() method with predefined message if `*err` == nil or
there's no match with `notErrs`.

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
func (l *Logger) CatchED(err *error, notErrs ...error) { l.catchED(true, err, notErrs...) }

func (l Logger) catchED(isDebugCatch bool, err *error, notErrs ...error) {
	if *err == nil {
		if isDebugCatch {
			l.D().Write(l.catchEDMsg)
		}

		return
	}

	ok := derr.IsInP(err, notErrs...)
	if !ok {
		l.E().Write((*err).Error())

		return
	}

	if isDebugCatch {
		l.D().Write(l.catchEDMsg)
	}
}
