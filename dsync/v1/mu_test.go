package dsync_test

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRWMutex(t *testing.T) {
	var (
		mu      = dsync.RWMutex{}
		group   = dsync.NewOneTimeGroup(dctx.New(dctx.OptionWithOSCancel()))
		n       = 0
		wantedN = 15000
	)

	for i := 0; i < wantedN/2; i++ {
		group.Go(func(context.Context) error {
			mu.LockF(func() {
				n++
			})

			return nil
		})

		group.Go(func(context.Context) error {
			mu.RLockF(func() {
				t.Log(n)
			})

			return nil
		})

		group.Go(func(context.Context) error {
			defer mu.Lock().Unlock()
			n++

			return nil
		})

		group.Go(func(context.Context) error {
			defer mu.RLock().RUnlock()
			t.Log(n)

			return nil
		})
	}

	err := group.Wait()
	require.NoError(t, err)
	require.EqualValues(t, wantedN, n)
}
