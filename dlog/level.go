package dlog

// Level - is a logs level pointing to respective Logger methods. Some levels can be set on/off in configuration.
type Level uint8

const (
	/*
		LevelError - it is used to log an error an operation met when this error cannot or shouldn't be returned.
		Otherwise, it's almost always better to return error and let an invoker decide what to do next compared to
		logging. This level is required and cannot be disabled with Logger configuration.
	*/
	LevelError Level = iota

	/*
		LevelWarn - it is used to log an error an operation met, and it's state was gracefully downgraded and execution
		process proceeded. This level is optional and may be disabled with logger configuration.
	*/
	LevelWarn

	/*
		LevelInfo - it is used to log a general purpose message such as build info, incoming requests metadata etc. This
		level is optional and may be disabled with logger configuration.
	*/
	LevelInfo

	/*
		LevelDebug - it is used to log an operation context in order to reveal execution details such as intermediate
		variables state. This level is optional and may be disabled with logger configuration.
	*/
	LevelDebug
)

// Is - reports whether `level` is completely enabled.
func (l Level) Is(level Level) bool {
	return (l & level) == 1
}

// Has - reports whether `level` is partly enabled.
func (l Level) Has(level Level) bool {
	return (l & level) > 0
}

// Enable - enables `level` provided.
func (l Level) Enable(level Level) Level {
	return l | level
}

// Disable - disables `level` provided. LevelError cannot be disabled.
func (l Level) Disable(level Level) Level {
	return l &^ level
}

func (l Level) String() string {
	switch l {
	case LevelError:
		return "E"

	case LevelWarn:
		return "W"

	case LevelInfo:
		return "I"

	case LevelDebug:
		return "D"

	default:
		return LevelError.String()
	}
}

func (l Level) Byte() byte {
	switch l {
	case LevelError:
		return 'E'

	case LevelWarn:
		return 'W'

	case LevelInfo:
		return 'I'

	case LevelDebug:
		return 'D'

	default:
		return LevelError.Byte()
	}
}
