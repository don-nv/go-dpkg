package dlog

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
)

const (
	// CatchEDDefaultMessage - is a default message used to be logged at LevelDebug for Logger.CatchED() method.
	CatchEDDefaultMessage = "OK"

	// TimeDefaultLayout - is the same as time.RFC3339Nano, but preserves trailing zeros.
	TimeDefaultLayout = "2006-01-02T15:04:05.000000000Z07:00"

	/*
		nameMaxBytes - is used as an empiric value to preallocate enough space while building a single Log name from
		names. See Log.constructName() method.

		Number is considered to be more or less accurate and is based on the following single name:
			- `handling_register_message_request`;
		- which is 33 bytes length. Give it 20% additional capacity and round up - this is the resulting value;

		TODO? This may be variable calculated over time for entire Logger or for each individual Log.
	*/
	nameMaxBytes = 40
)

// ReadScopeFn - is used at Log.Scope() method.
type ReadScopeFn func(ctx context.Context, log Logger) Logger

/*
ReadScopeDefault - default ReadScopeFn function. Uses dctx package to populate Logger with:
  - Goroutine id;
  - X request id;

- both values are expected to be set to `ctx` via dctx package. If some of them are missing, then the respective empty
values get omitted.
*/
func ReadScopeDefault(ctx context.Context, log Logger) Logger {
	var data = log.With()

	id := dctx.GoID(ctx)
	if id != "" {
		data = data.String("go_id", id)
	}

	id = dctx.XRequestID(ctx)
	if id != "" {
		data = data.String("x_req_id", id)
	}

	return data.Build()
}

// errStr - returns error message if not nil, otherwise - "nil".
func errStr(err error) string {
	if err != nil {
		return err.Error()
	}

	return "nil"
}
