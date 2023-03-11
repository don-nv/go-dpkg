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

func TestRoundRobinSync_Next(t *testing.T) {
	var (
		ctx = dctx.New(dctx.OptionWithOSCancel())

		roundsN = 10000
		robin   = dstruct.NewRoundRobinSync(
			dstruct.RoundRobin{
				RoundsMaxN: roundsN,
			},
		)
	)

	var n int
	// Two subsequent loops.
	for i := 0; i < 2; i++ {
		for j := 0; j < roundsN; j++ {
			require.EqualValues(t, n, robin.NextI())
			n++

			if n == roundsN {
				n = 0
			}
		}
	}

	var (
		group   = dsync.NewGroup(ctx)
		roundIC = make(chan int, roundsN)
	)
	go func() {
		time.Sleep(time.Second) // Give some time to launch goroutines.

		require.NoError(t, group.Wait())
		close(roundIC)
	}()

	for i := 0; i < roundsN; i++ {
		group.Go(func(context.Context) error {
			roundIC <- robin.NextI()

			return nil
		})
		group.Go(func(context.Context) error {
			roundIC <- robin.NextI()

			return nil
		})
	}

	var roundsSum int

	for i := 0; i < roundsN; i++ {
		roundsSum += i
	}
	roundsSum *= 2 // Two goroutines were sending to the channel.

	for i := range roundIC {
		roundsSum -= i
	}
	require.EqualValues(t, 0, roundsSum)
}
