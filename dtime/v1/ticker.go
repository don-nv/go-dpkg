package dtime

import (
	"context"
	"github.com/don-nv/go-dpkg/dchan/v1"
	"time"
)

type Ticker struct {
	ticker *time.Ticker
}

func NewTicker(d time.Duration) Ticker {
	return Ticker{
		ticker: time.NewTicker(d),
	}
}

func (t Ticker) C() <-chan time.Time {
	return t.ticker.C
}

// Stop - is the same as time.Ticker.Stop(), but drains tickers' channel after.
func (t Ticker) Stop() {
	t.ticker.Stop()

	dchan.Drain(context.Background(), t.ticker.C)
}

// Reset - is the same as time.Ticker.Reset().
func (t Ticker) Reset(d time.Duration) {
	t.ticker.Reset(d)
}
