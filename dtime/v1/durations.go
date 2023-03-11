package dtime

import (
	"context"
	"fmt"
	"time"
)

func SinceMs[T float64](t time.Time) T {
	return T(time.Since(t).Milliseconds())
}

// AwaitDelay - awaits 'd'. If 'ctx' is canceled before 'd', context error is returned.
func AwaitDelay(ctx context.Context, d time.Duration) error {
	if d < 1 {
		return fmt.Errorf("non-positive %q delay", d)
	}

	var ticker = NewTicker(d)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-ticker.C():
		return nil
	}
}
