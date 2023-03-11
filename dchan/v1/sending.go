package dchan

import (
	"context"
	"sync"
)

// TrySend - tries to send T to 'c' without blocking. Returned bool indicates if sending succeeded.
func TrySend[T any](ctx context.Context, c chan<- T, t T) bool {
	select {
	case <-ctx.Done():
		return false

	case c <- t:
		return true

	default:
		return false
	}
}

// Send - sends T to 'c' until 'ctx' done. Returned bool indicates if sending succeeded.
func Send[T any](ctx context.Context, c chan<- T, t T) bool {
	select {
	case <-ctx.Done():
		return false

	case c <- t:
		return true
	}
}

/*
FanningOut - receives from 'c' and sends to one of 'cs' until 'ctx' is not done. 'cs' not get closed after. Each T
received from 'c' is sent to one of 'cs' only.
*/
func FanningOut[T any](ctx context.Context, c <-chan T, cs ...chan<- T) {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(len(cs))
	for _, outC := range cs {
		var outC = outC

		go func() {
			defer wg.Done()

			for {
				t, ok := Receive(ctx, c)
				if !ok {
					return
				}

				ok = Send(ctx, outC, t)
				if !ok {
					return
				}
			}
		}()
	}
}
