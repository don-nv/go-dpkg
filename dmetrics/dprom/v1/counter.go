package dprom

import "github.com/prometheus/client_golang/prometheus"

var counterVec *prometheus.CounterVec

func init() {
	counterVec = newCounterVec("counter", labelsSlice)

	prometheus.MustRegister(counterVec)

	// Check robustness. May panic at runtime.
	NewEntry.Counter().Add(-1)
}

type Counter struct {
	entry Entry
}

func newCounter(entry Entry) Counter {
	return Counter{
		entry: entry,
	}
}

func (c Counter) Inc() {
	c.entry.promLabelsTmp(
		func(labels prometheus.Labels) {
			counterVec.
				With(labels).
				Inc()
		},
	)
}

func (c Counter) Add(v float64) {
	// May panic otherwise.
	if v < 0 {
		v = 0
	}

	c.entry.promLabelsTmp(
		func(labels prometheus.Labels) {
			counterVec.
				With(labels).
				Add(v)
		},
	)
}
