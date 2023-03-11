package roundbreaker

import (
	"context"
	"github.com/don-nv/go-dpkg/dstruct/v1"
	"runtime"
)

// Rounds - is a Round-Robin for hosts.
type Rounds struct {
	hosts []Host
	robin dstruct.RoundRobinSync
}

func NewRounds(hosts []Host) *Rounds {
	return &Rounds{
		hosts: hosts,
		robin: dstruct.NewRoundRobinSync(dstruct.RoundRobin{
			RoundsMaxN: len(hosts),
		}),
	}
}

// nextHost - returns next enabled host in rounds or context error.
func (r *Rounds) nextHost(ctx context.Context) (*Host, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		default:
			var i = r.robin.NextI()
			if r.hosts[i].enabled() {
				return &r.hosts[i], nil
			}

			runtime.Gosched()
		}
	}
}
