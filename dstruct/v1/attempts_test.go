package dstruct_test

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dstruct/v1"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//nolint:funlen
func TestAttemptsSync(t *testing.T) {
	var (
		ctx = dctx.New(dctx.OptionWithOSCancel())

		delays = []time.Duration{
			100 * time.Millisecond,
			200 * time.Millisecond,
			300 * time.Millisecond,
		}
		attempts = dstruct.NewAttemptsV1Sync(dstruct.AttemptsV1{
			Delays: delays,
		})
	)

	// Second run performs after attempts reset.
	for i := 0; i < 2; i++ {
		var ok = attempts.Next()
		require.True(t, ok)

		var n = 1
		for ok {
			require.EqualValues(t, n, attempts.AttemptN())
			require.EqualValues(t, len(delays), attempts.AttemptsN())
			require.EqualValues(t, delays[n-1], attempts.Delay())

			var startedAt = time.Now()

			err := attempts.AwaitDelay(ctx)
			require.NoError(t, err)

			var since = time.Since(startedAt)
			require.Greater(t, since, attempts.Delay())
			require.Less(t, since, attempts.Delay()+2*time.Millisecond)

			ok = attempts.Next()
			n++
		}

		attempts.Reset()
	}

	var (
		group     = dsync.NewGroup(ctx)
		startedAt = time.Now()
	)

	for i := 0; i < 2; i++ {
		for attempts.Next() {
			var c = make(chan struct{})

			group.Go(func(context.Context) error {
				close(c)

				return attempts.AwaitDelay(ctx)
			})

			<-c
		}

		attempts.Reset()
	}

	require.NoError(t, group.Wait())

	var (
		delayMax = delays[len(delays)-1]
		since    = time.Since(startedAt)
	)
	require.Greater(t, since, delayMax)
	require.Less(t, since, delayMax+4*time.Millisecond)
}
