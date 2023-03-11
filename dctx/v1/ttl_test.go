package dctx_test

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWithTTLLessTimeout(t *testing.T) {
	t.Parallel()

	data := []struct {
		//nolint:containedctx
		Context      context.Context
		TimeoutToSet *time.Duration
	}{
		{
			Context: dctx.New(),
		},
		{
			Context:      dctx.New(),
			TimeoutToSet: func() *time.Duration { var d = time.Second; return &d }(),
		},
		{
			Context: func() context.Context {
				var ctx, _ = dctx.NewTTL(dctx.OptionTTLWithTimeout(2 * time.Second))

				return ctx
			}(),
			TimeoutToSet: func() *time.Duration { var d = time.Second; return &d }(),
		},
	}

	for _, d := range data {
		var d = d

		t.Run("", func(t *testing.T) {
			t.Parallel()

			var (
				ctx            = d.Context
				cancel         func()
				deadlineWanted = d.TimeoutToSet != nil
			)

			if deadlineWanted {
				ctx, cancel = dctx.WithTTLTimeout(d.Context, *d.TimeoutToSet)
				require.NotNil(t, cancel)

				defer cancel()
			}

			_, ok := ctx.Deadline()
			require.EqualValues(t, deadlineWanted, ok)

			var ttl = dctx.Timeout(ctx)

			if deadlineWanted {
				require.Less(t, ttl, *d.TimeoutToSet)
				require.Greater(t, ttl, *d.TimeoutToSet-time.Millisecond)
			} else {
				require.Zero(t, ttl)
			}
		})
	}
}
