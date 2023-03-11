package dprom

import (
	"github.com/don-nv/go-dpkg/dtime/v1"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var histogramMsVec *prometheus.HistogramVec

func init() {
	histogramMsVec = newHistogramMsVec("histogram_ms", labelsSlice)

	prometheus.MustRegister(histogramMsVec)

	// Check robustness. May panic at runtime.
	NewEntry.HistogramMs().Write()
}

type HistogramMs struct {
	entry     Entry
	startedAt time.Time
}

func newHistogramMs(entry Entry) HistogramMs {
	return HistogramMs{
		entry:     entry,
		startedAt: time.Now(),
	}
}

func (h HistogramMs) Write() {
	var d = dtime.SinceMs(h.startedAt)

	h.entry.promLabelsTmp(
		func(labels prometheus.Labels) {
			histogramMsVec.
				With(labels).
				Observe(d)
		},
	)
}

func (h HistogramMs) WriteResult(v string) {
	var d = dtime.SinceMs(h.startedAt)

	h.entry.
		WithResult(v).
		promLabelsTmp(
			func(labels prometheus.Labels) {
				histogramMsVec.
					With(labels).
					Observe(d)
			},
		)
}
