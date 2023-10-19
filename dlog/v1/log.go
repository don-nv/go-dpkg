package dlog

import (
	"context"
	"github.com/rs/zerolog"
	"strings"
)

type Log struct {
	writeLvl Level
	logger   Logger
	data     Data
}

func newLog(log Logger, writeLvl Level) Log {
	return Log{
		writeLvl: writeLvl,
		logger:   log,
		data:     newData(log),
	}
}

func (l Log) Write(msg string) Log {
	if l.logger.zero.GetLevel() == zerolog.Disabled {
		return l
	}

	newEvent, ok := eventFactoryByLvl[l.writeLvl]
	if !ok {
		newEvent = eventFactoryByLvl[LevelError]
	}

	var event = newEvent(l.logger.zero)

	name := l.constructName()
	if name != "" {
		event.Str("name", name)
	}

	event.Timestamp().Msg(msg)

	return l
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

func (l Log) constructName() string {
	if len(l.logger.names) < 1 {
		return ""
	}

	// TODO? Grow may be replaced with bytes sync pools.

	var builder = strings.Builder{}
	builder.Grow(len(l.logger.names) * nameMaxBytes)

	for i, name := range l.logger.names {
		builder.WriteString(name)

		if i != len(l.logger.names)-1 {
			builder.WriteByte('.')
		}
	}

	return builder.String()
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
