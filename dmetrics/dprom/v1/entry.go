package dprom

import "github.com/prometheus/client_golang/prometheus"

var NewEntry Entry

type Entry struct {
	entityGroup string
	entity      string
	methodGroup string
	method      string
	result      string
}

func (e Entry) WithEntityGroup(v string) Entry {
	e.entityGroup = v

	return e
}

func (e Entry) WithEntity(v string) Entry {
	e.entity = v

	return e
}

func (e Entry) WithMethodGroup(v string) Entry {
	e.methodGroup = v

	return e
}

func (e Entry) WithMethod(v string) Entry {
	e.method = v

	return e
}

func (e Entry) WithResult(v string) Entry {
	e.result = v

	return e
}

// promLabelsTmp - provides prometheus labels in 'f'. Provided labels must not be used after 'f' returns.
func (e Entry) promLabelsTmp(f func(labels prometheus.Labels)) {
	var labels = promLabelsPool.Get()
	defer promLabelsPool.Put(labels)

	labels[labelEntityGroup] = e.entityGroup
	labels[labelEntity] = e.entity
	labels[labelMethodGroup] = e.methodGroup
	labels[labelMethod] = e.method
	labels[labelResult] = e.result

	f(labels)
}

func (e Entry) Counter() Counter {
	return newCounter(e)
}

func (e Entry) HistogramMs() HistogramMs {
	return newHistogramMs(e)
}

func (e Entry) Gauge() Gauge {
	return newGauge(e)
}
