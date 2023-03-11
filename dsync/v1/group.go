package dsync

import (
	"context"
	"fmt"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/derr/v1"
	"golang.org/x/sync/errgroup"
)

// Group - is a wrap-around errgroup.WithContext providing additional methods.
type Group struct {
	group *errgroup.Group
	doneC <-chan struct{}
}

// NewGroup - returns a new Group. The Group is Done() once 'ctx' is closed. May be used several times.
func NewGroup(ctx context.Context) Group {
	return Group{
		group: &errgroup.Group{},
		doneC: ctx.Done(),
	}
}

/*
NewOneTimeGroup - returns a new Group. The Group is permanently Done() the first time a function passed to Go() (etc)
returns a non-nil error, 'ctx' is canceled or the first time Wait() returns, whichever occurs first.
*/
func NewOneTimeGroup(ctx context.Context) Group {
	var group, groupCtx = errgroup.WithContext(ctx)

	return Group{
		group: group,
		doneC: groupCtx.Done(),
	}
}

/*
WaitE - is the same as Wait(), but is intended to be used with 'defer' capturing returned 'err'. If Wait() returns err
!= nil, then this error is joined to 'err'.
*/
func (g Group) WaitE(err *error) {
	werr := g.group.Wait()
	derr.JoinInP(err, werr)
}

// CatchE - is the same as WaitE(), but if captured 'err' != nil, then group gets canceled.
func (g Group) CatchE(err *error) {
	if !derr.InP(err, nil) {
		g.Cancel()
	}

	werr := g.group.Wait()
	derr.JoinInP(err, werr)
}

func (g Group) Wait() error {
	return g.group.Wait()
}

func (g Group) Cancel() {
	g.group.Go(func() error {
		return context.Canceled
	})
}

func (g Group) CancelE(err error) {
	if err == nil {
		err = context.Canceled
	}

	g.group.Go(func() error {
		return err
	})
}

func (g Group) Done() bool {
	select {
	case _, ok := <-g.doneC:
		return !ok

	default:
		return false
	}
}

func (g Group) Go(f func(ctx context.Context) error) {
	if g.Done() {
		g.CancelE(derr.ErrDone)
	}

	// Use group context.
	var ctx, cancel = g.DeriveContext()

	g.group.Go(func() error {
		defer cancel()

		return f(ctx)
	})
}

func (g Group) GoTry(f func(ctx context.Context) error) bool {
	if g.Done() {
		g.CancelE(derr.ErrDone)
	}

	var ctx, cancel = g.DeriveContext()

	return g.group.TryGo(func() error {
		defer cancel()

		return f(ctx)
	})
}

func (g Group) SetLimit(n int) {
	g.group.SetLimit(n)
}

/*
GoUntilWait - is the same as Go(), but if 'f' exits without error before Wait(), a whole group is cancelled and an
error is returned in Wait(). This method is intended to be used for background goroutines that must last until group is
not done.
  - 'name' - is used to distinct 'f's from each other, that is gets written into errors value;
*/
func (g Group) GoUntilWait(name string, f func(ctx context.Context) error) {
	var ctx, cancel = g.DeriveContext()

	g.group.Go(func() error {
		defer cancel()

		err := f(ctx)
		if err != nil {
			return err
		}

		if ctx.Err() == nil {
			return fmt.Errorf("%q has stopped before group wait", name)
		}

		return nil
	})
}

/*
DeriveContext - returns new context that will be cancelled if Group exits. A special care should be taken not to
forget to release new context by calling returned cancel function if Group lifetime is long. Context deriving without
control may lead to goroutines leaking.
*/
func (g Group) DeriveContext() (context.Context, context.CancelFunc) {
	var ctx, cancel = dctx.NewTTL()

	ctx, _ = dctx.WithTTLC(ctx, g.doneC)

	return ctx, cancel
}
