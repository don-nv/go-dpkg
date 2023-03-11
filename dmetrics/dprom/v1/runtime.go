package dprom

import (
	"github.com/don-nv/go-dpkg/dctx/v1"
	"runtime"
	"time"
)

//nolint:funlen
func init() {
	if !configEnv.GORuntime.Enabled {
		return
	}

	var (
		ctx, cancel = dctx.NewTTL(dctx.OptionTTLWithOSCancel())
		entry       = NewEntry.WithEntityGroup("go_runtime")
	)

	go func() {
		var (
			goroutinesLastC = make(chan float64, 2)

			heapBytesAllocatedC   = make(chan float64, 2)
			heapBytesRetainedC    = make(chan float64, 2)
			heapBytesReservedC    = make(chan float64, 2)
			heapObjectsAllocatedC = make(chan float64, 2)
			heapObjectsAliveC     = make(chan float64, 2)

			gcBytesUntilNextGCC = make(chan float64, 2)
			gcCPUConsumptionC   = make(chan float64, 2)
			gcCyclesLastC       = make(chan float64, 2)
			gcPauseLastNsC      = make(chan float64, 2)
		)

		var goroutines = entry.WithEntity("goroutines")
		goroutines.WithMethod("last").Gauge().RunSettingC(goroutinesLastC)

		var heap = entry.WithEntity("heap")
		heap.WithMethod("bytes_allocated").Gauge().RunSettingC(heapBytesAllocatedC)
		heap.WithMethod("bytes_retained").Gauge().RunSettingC(heapBytesRetainedC)
		heap.WithMethod("bytes_reserved").Gauge().RunSettingC(heapBytesReservedC)
		heap.WithMethod("objects_alive").Gauge().RunSettingC(heapObjectsAliveC)
		heap.WithMethod("objects_allocated").Gauge().RunSettingC(heapObjectsAllocatedC)

		var gc = entry.WithEntity("gc")
		gc.WithMethod("bytes_until_next_gc").Gauge().RunSettingC(gcBytesUntilNextGCC)
		gc.WithMethod("cpu_consumption").Gauge().RunSettingC(gcCPUConsumptionC)
		gc.WithMethod("cycles_last").Gauge().RunSettingC(gcCyclesLastC)
		gc.WithMethod("pause_last_ns").Gauge().RunSettingC(gcPauseLastNsC)

		defer func() {
			<-ctx.Done()
			cancel()

			close(goroutinesLastC)

			close(heapBytesAllocatedC)
			close(heapBytesRetainedC)
			close(heapBytesReservedC)
			close(heapObjectsAliveC)
			close(heapObjectsAllocatedC)

			close(gcCPUConsumptionC)
			close(gcBytesUntilNextGCC)
			close(gcCyclesLastC)
			close(gcPauseLastNsC)
		}()

		for {
			var stats = NewGoRuntime()

			select {
			case <-ctx.Done():
				return

			case <-time.After(configEnv.GORuntime.Delay):
				stats = stats.Refresh()

				goroutinesLastC <- stats.GoroutinesLast

				heapBytesAllocatedC <- stats.HeapBytesAllocated
				heapBytesRetainedC <- stats.HeapBytesRetained
				heapBytesReservedC <- stats.HeapBytesReserved
				heapObjectsAliveC <- stats.HeapObjectsAlive
				heapObjectsAllocatedC <- stats.HeapObjectsAllocated

				gcBytesUntilNextGCC <- stats.GCBytesUntilNextGC
				gcCPUConsumptionC <- stats.GCCPUConsumption
				gcCyclesLastC <- stats.GCCyclesLast
				gcPauseLastNsC <- stats.GCPauseLastNs
			}
		}
	}()
}

type GoRuntime struct {
	GoroutinesLast float64

	HeapBytesAllocated   float64
	HeapBytesRetained    float64
	HeapBytesReserved    float64
	HeapObjectsAlive     float64
	HeapObjectsAllocated float64

	GCLastFinishedAt   float64
	GCBytesUntilNextGC float64
	GCCPUConsumption   float64
	GCCyclesTotal      float64
	GCCyclesLast       float64
	GCPauseTotalNs     float64
	GCPauseLastNs      float64
}

func NewGoRuntime() GoRuntime {
	var stats = runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	return GoRuntime{
		GoroutinesLast: float64(runtime.NumGoroutine()),

		HeapBytesAllocated:   float64(stats.HeapAlloc),
		HeapBytesRetained:    float64(stats.HeapIdle - stats.HeapReleased),
		HeapBytesReserved:    float64(stats.HeapSys),
		HeapObjectsAlive:     float64(stats.Mallocs - stats.Frees),
		HeapObjectsAllocated: float64(stats.HeapObjects),

		GCLastFinishedAt:   float64(stats.LastGC),
		GCBytesUntilNextGC: float64(stats.NextGC),
		GCCPUConsumption:   stats.GCCPUFraction,
		GCCyclesTotal:      float64(stats.NumGC),
		GCCyclesLast:       float64(stats.NumGC),
		GCPauseTotalNs:     float64(stats.PauseTotalNs),
		GCPauseLastNs:      float64(stats.PauseTotalNs),
	}
}

func (r GoRuntime) Refresh() GoRuntime {
	var stats = runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	r.GoroutinesLast = float64(runtime.NumGoroutine())

	r.HeapBytesAllocated = float64(stats.HeapAlloc)
	r.HeapBytesRetained = float64(stats.HeapIdle - stats.HeapReleased)
	r.HeapBytesReserved = float64(stats.HeapSys)
	r.HeapObjectsAlive = float64(stats.Mallocs - stats.Frees)
	r.HeapObjectsAllocated = float64(stats.HeapObjects)

	r.GCLastFinishedAt = float64(stats.LastGC)
	r.GCBytesUntilNextGC = float64(stats.NextGC)
	r.GCCPUConsumption = stats.GCCPUFraction

	r.GCCyclesLast = float64(stats.NumGC) - r.GCCyclesTotal
	r.GCCyclesTotal = float64(stats.NumGC)

	r.GCPauseLastNs = float64(stats.PauseTotalNs) - r.GCPauseTotalNs
	r.GCPauseTotalNs = float64(stats.PauseTotalNs)

	return r
}
