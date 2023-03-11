package dchan

import (
	"context"
	"sync"
)

// Drain - drains 'c' until there are Ts in 'c' or 'ctx' done. Returned bool indicates if draining succeeded.
func Drain[T any](ctx context.Context, c <-chan T) bool {
	for {
		select {
		case <-ctx.Done():
			return false

		case _, ok := <-c:
			if !ok {
				return true
			}

		default:
			return true
		}
	}
}

/*
Closed - reports whether 'c' is closed or not. If 'c' is open, then first element in the chan may become last. Also,
this method may block current goroutine until 'c' is not able to receive read first element.
*/
func Closed[T any](c chan T) bool {
	select {
	case t, ok := <-c:
		if ok {
			Send(context.Background(), c, t)
		}

		return !ok

	default:

		return false
	}
}

// TryReceive - tries to receive T from 'c' without blocking. Returned bool indicates if receiving succeeded.
func TryReceive[T any](ctx context.Context, c <-chan T) (T, bool) {
	select {
	case <-ctx.Done():
		var t T
		return t, false

	case t, ok := <-c:
		return t, ok

	default:
		var t T
		return t, false
	}
}

// Receive - receives T from 'c' until 'ctx' done or 'c' closed. Returned bool indicates if receiving succeeded.
func Receive[T any](ctx context.Context, c <-chan T) (T, bool) {
	select {
	case <-ctx.Done():
		var t T
		return t, false

	case t, ok := <-c:
		return t, ok
	}
}

// NewFanIn - is the same as FanningIn, but creates chan and returns it without block.
func NewFanIn[T any](ctx context.Context, buff int, cs ...<-chan T) chan<- T {
	if buff < 0 {
		buff = 1
	}

	var c = make(chan T, buff)

	go FanningIn[T](ctx, c, cs...)

	return c
}

/*
FanningIn - receives from 'cs' and sends to 'c' until 'ctx' is not done or there are no open channels in 'cs'. 'c' is
not closed after and must not be closed until function returns.
*/
func FanningIn[T any](ctx context.Context, c chan<- T, cs ...<-chan T) {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(len(cs))
	for _, inC := range cs {
		var inC = inC

		go func() {
			defer wg.Done()

			for {
				t, ok := Receive(ctx, inC)
				if !ok {
					return
				}

				ok = Send(ctx, c, t)
				if !ok {
					return
				}
			}
		}()
	}
}
