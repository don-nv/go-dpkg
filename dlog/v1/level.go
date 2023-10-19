package dlog

type Level uint8

/*
LevelError - is used to log an error an operation met. An error should be logged when it cannot or shouldn't be
returned. Otherwise, it's almost always better to return error and let a caller decide what to do next.
*/
const LevelError Level = 0

const (
	/*
		LevelWarn - is used to log an error (attention message) an operation met. Error here may be not a regular error,
		but an unexpected operation state. This state can be downgraded gracefully and this operation may proceed.
	*/
	LevelWarn Level = 1 << iota

	// LevelInfo - is used to log a general purpose message: build info, incoming requests metadata etc.
	LevelInfo

	// LevelDebug - is used to log an operation context revealing execution details: variables, operations etc.
	LevelDebug
)

// LevelAll - represents all enabled levels.
const LevelAll = LevelError | LevelWarn | LevelInfo | LevelDebug

// NewLevel - is the same as Level.Enable(), but provides verbose way of Level initialization.
func NewLevel(levels ...Level) Level {
	return LevelError.Enable(levels...)
}

// Enabled - reports whether `levels` are enabled or not.
func (l Level) Enabled(levels ...Level) bool {
	for _, level := range levels {
		if level == LevelError { // enabled by default.
			continue
		}

		if l&level == level {
			continue
		}

		return false
	}

	return true
}

// Enable - enables `levels` provided. LevelError is enabled by default.
func (l Level) Enable(levels ...Level) Level {
	for _, lvl := range levels {
		l |= lvl
	}

	return l
}

// Disable - disables `levels` provided. LevelError cannot be disabled.
func (l Level) Disable(levels ...Level) Level {
	for _, level := range levels {
		l &^= level
	}

	return l
}

// String - returns Level string representation. If Level is unknown, a value for LevelError is returned.
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
