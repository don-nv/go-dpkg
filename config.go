package dpkg

import (
	"os"
)

// EnvDebugEnabled - is a key of env variable. If env variable contains 'true', then debug is enabled.
const EnvDebugEnabled = "DEBUG_ENABLED"

var debugEnabled = func() bool {
	v, ok := os.LookupEnv(EnvDebugEnabled)
	if !ok {
		return true
	}

	return v == "true"
}()

func DebugEnabled() bool {
	return debugEnabled
}
