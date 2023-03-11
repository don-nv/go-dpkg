package dstruct

import (
	"github.com/don-nv/go-dpkg/dsync/v1"
)

// RoundRobinSync - is the same as RoundRobin, but it is safe to be used concurrently.
type RoundRobinSync struct {
	robin RoundRobin
	mu    dsync.RWMutex
}

func NewRoundRobinSync(robin RoundRobin) RoundRobinSync {
	return RoundRobinSync{
		robin: robin,
	}
}

// NextI - returns next round index.
func (r *RoundRobinSync) NextI() int {
	defer r.mu.Lock().Unlock()

	return r.robin.NextI()
}

type RoundRobin struct {
	i          int
	RoundsMaxN int
}

// NextI - returns next round index.
func (r *RoundRobin) NextI() int {
	if r.i > r.RoundsMaxN-1 {
		r.i = 0
	}

	var i = r.i

	r.i++

	return i
}
