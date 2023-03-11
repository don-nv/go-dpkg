package dlog

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
)

const (
	HiddenValueString = "?"
	HiddenValueByte   = '?'

	// CatchEDDefaultMessage - is a default message used to be logged at LevelDebug for Logger.CatchED() method.
	CatchEDDefaultMessage = "OK"

	// TimeDefaultLayout - is the same as time.RFC3339Nano, but preserves trailing nanoseconds zeros.
	TimeDefaultLayout = "2006-01-02T15:04:05.000000000Z07:00"

	/*
		nameExpectedMaxBytes - is used as an empiric value to preallocate enough space while building a single Log name
		from names. See Log.constructName() method.

		Number is considered to be more or less accurate and is based on the following single name:
			- 'handling_register_message_request';
		It has length of 33 bytes. Give it 20% additional capacity and round up - this is the resulting value;

	*/
	nameExpectedMaxBytes = 40
	// TODO? This may be variable calculated over time for entire Logger or for each individual Log.
)

// ReadScopeFn - is used at Data.Scope() or Log.Scope() method.
type ReadScopeFn func(ctx context.Context, data Data) Data

/*
ReadScopeDefault - default ReadScopeFn function. Uses dctx package to populate Logger with:
  - Goroutine id;
  - X request id;

Both values are expected to be set to 'ctx' via dctx package. If some of them are missing, then the respective empty
values get omitted.
*/
func ReadScopeDefault(ctx context.Context, data Data) Data {
	id := dctx.GoID(ctx)
	if id != "" {
		data = data.String("go_id", id)
	}

	id = dctx.XRequestID(ctx)
	if id != "" {
		data = data.String("x_req_id", id)
	}

	deadline, ok := ctx.Deadline()
	if ok {
		data = data.String("deadline", deadline.String())
	}

	return data
}
