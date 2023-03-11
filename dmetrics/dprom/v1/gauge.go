package dprom

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var gaugeVec *prometheus.GaugeVec

func init() {
	gaugeVec = newGaugeVec("gauge", labelsSlice)

	prometheus.MustRegister(gaugeVec)

	// Check robustness. May panic at runtime.
	NewEntry.Gauge().Set(0)
}

type Gauge struct {
	entry Entry
}

func newGauge(entry Entry) Gauge {
	return Gauge{
		entry: entry,
	}
}

func (g Gauge) Inc() {
	g.entry.promLabelsTmp(
		func(labels prometheus.Labels) {
			gaugeVec.
				With(labels).
				Inc()
		},
	)
}

func (g Gauge) Dec() {
	g.entry.promLabelsTmp(
		func(labels prometheus.Labels) {
			gaugeVec.
				With(labels).
				Dec()
		},
	)
}

func (g Gauge) Add(v float64) {
	g.entry.promLabelsTmp(
		func(labels prometheus.Labels) {
			gaugeVec.
				With(labels).
				Add(v)
		},
	)
}

func (g Gauge) Set(v float64) {
	g.entry.promLabelsTmp(
		func(labels prometheus.Labels) {
			gaugeVec.
				With(labels).
				Set(v)
		},
	)
}

func (g Gauge) RunSettingC(c <-chan float64) {
	go g.SettingC(c)
}

func (g Gauge) SettingC(c <-chan float64) {
	for v := range c {
		g.Set(v)
	}
}

func (g Gauge) SettingF(ctx context.Context, f func(ctx context.Context) float64, delay time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(delay):
			g.Set(f(ctx))
		}
	}
}
