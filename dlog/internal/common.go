package internal

import (
	"context"
	"dpkg/derr"
	"dpkg/dlog"
)

func ReadCtx(ctx context.Context, log dlog.Logger, readCtx ReadContext) dlog.Logger {
	for _, kv := range readCtx(ctx) {
		log.WithKV(kv.Key, kv.Value)
	}

	return log
}

/*
CatchE - logs error at the respective level if `*err != nil`. Otherwise, if `useDebugLevel`, logs predefined message
at debug level.
*/
func CatchE(log dlog.Logger, err *error, useDebugLevel bool, expectedErrors ...error) dlog.Logger {
	var (
		caught bool
		errVal = derr.Dereference(err)
	)

	defer func() {
		if caught {
			log.E(errVal)

			return
		}

		if useDebugLevel {
			log.D("(context)")

			return
		}
	}()

	if errVal == nil {
		return log
	}

	caught = !derr.Is(errVal, expectedErrors...)

	return log
}

func Flushing(ctx context.Context, log dlog.Logger, sync func(ctx context.Context) error) dlog.Logger {
	defer log.WithContext(ctx).WithName("logger", "flushing").I("running...").I("...done")

	var doneC = make(chan struct{})

	go func() { // TODO: Replace with dgo.
		defer close(doneC)

		err := sync(ctx)
		if err != nil {
			//log.E("syncing buffered log entries: %s", err)
		}
	}()

	select {
	case _ = <-ctx.Done():
		select {
		case <-doneC:

		default:
			//log.E("incomplete: %s", err)
		}

	case <-doneC:
	}

	return log
}
