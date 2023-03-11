package dprom

import (
	"github.com/don-nv/go-dpkg/dsync/v1"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelEntityGroup = "entity_group"
	labelEntity      = "entity"
	labelMethodGroup = "method_group"
	labelMethod      = "method"
	labelResult      = "result"
)

var labelsSlice = []string{
	labelEntityGroup,
	labelEntity,
	labelMethodGroup,
	labelMethod,
	labelResult,
}

var promLabelsPool = dsync.NewPool[prometheus.Labels](
	func() prometheus.Labels {
		var labels = make(prometheus.Labels, len(labelsSlice))

		for _, l := range labelsSlice {
			labels[l] = ""
		}

		return labels
	},
	dsync.OptionPoolWithOnPut(func(labels prometheus.Labels) {
		for k := range labels {
			labels[k] = ""
		}
	}),
)
