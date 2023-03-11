package dprom

import "github.com/prometheus/client_golang/prometheus"

const vectorNameSuffix = "_dprom"

func newCounterVec(name string, labels []string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: configEnv.Namespace,
			Subsystem: configEnv.Subsystem,
			Name:      name + vectorNameSuffix,
		},
		labels,
	)
}

func newHistogramMsVec(name string, labels []string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		newHistOptions(name, configEnv.BucketsPartitionsMs),
		labels,
	)
}

func newHistOptions(name string, buckets []float64) prometheus.HistogramOpts {
	return prometheus.HistogramOpts{
		Namespace: configEnv.Namespace,
		Subsystem: configEnv.Subsystem,
		Name:      name + vectorNameSuffix,
		Buckets:   buckets,
	}
}

func newGaugeVec(name string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: configEnv.Namespace,
			Subsystem: configEnv.Subsystem,
			Name:      name + vectorNameSuffix,
		},
		labels,
	)
}
