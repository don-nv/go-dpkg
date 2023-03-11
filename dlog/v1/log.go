package dlog

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type Log struct {
	startedAt time.Time
	writeLvl  Level
	data      Data
}

func E() Log {
	return New().E()
}

func newLog(log Logger, writeLvl Level) Log {
	return Log{
		writeLvl:  writeLvl,
		data:      newData(log),
		startedAt: log.startedAt,
	}
}

func (l Log) Stack() Log {
	l.data = l.data.Stack()

	return l
}

func (l Log) Writef(format string, args ...interface{}) Log {
	return l.Write(fmt.Sprintf(format, args...))
}

func (l Log) Write(msg string) Log {
	var (
		logger = l.data.Build()
		event  = l.newEventFactoryForLevel()(logger.zero)
	)

	name := logger.constructName()
	if name != "" {
		event = event.Str("name", name)
	}

	if !l.startedAt.IsZero() {
		event = event.Str("duration", time.Since(l.startedAt).String())
	}

	event.Timestamp().Msg(msg)

	return l
}

func (l Log) newEventFactoryForLevel() func(l zerolog.Logger) *zerolog.Event {
	newEvent, ok := eventFactoryByLvl[l.writeLvl]
	if !ok {
		newEvent = eventFactoryByLvl[LevelError]
	}

	return newEvent
}

var eventFactoryByLvl = map[Level]func(l zerolog.Logger) *zerolog.Event{
	LevelError: func(l zerolog.Logger) *zerolog.Event {
		return l.Error()
	},
	LevelWarn: func(l zerolog.Logger) *zerolog.Event {
		return l.Warn()
	},
	LevelInfo: func(l zerolog.Logger) *zerolog.Event {
		return l.Info()
	},
	LevelDebug: func(l zerolog.Logger) *zerolog.Event {
		return l.Debug()
	},
}

// Scope - is the same as Data.Scope().
func (l Log) Scope(ctx context.Context) Log {
	l.data = l.data.Scope(ctx)

	return l
}

// Name - is the same as Data.Name().
func (l Log) Name(names ...string) Log {
	l.data = l.data.Name(names...)

	return l
}

// Any - is the same as Data.Any().
func (l Log) Any(key string, value interface{}) Log {
	l.data = l.data.Any(key, value)

	return l
}
