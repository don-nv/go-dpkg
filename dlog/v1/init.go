package dlog

import "github.com/rs/zerolog"

func init() {
	zerolog.LevelErrorValue = LevelError.String()
	zerolog.LevelWarnValue = LevelWarn.String()
	zerolog.LevelInfoValue = LevelInfo.String()
	zerolog.LevelDebugValue = LevelDebug.String()
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = TimeDefaultLayout
}
