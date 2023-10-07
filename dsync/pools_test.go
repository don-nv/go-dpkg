package dsync_test

import (
	"github.com/don-nv/go-dpkg/dsync"
	"testing"
)

func BenchmarkNewBytesPools(b *testing.B) {
	defer b.ReportAllocs()

	var pools = dsync.NewPoolsBytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 10; i++ {
			var bytes = pools.Acquire(i)

			//fmt.Println("id: ", i, "", cap(bytes))
			bytes = append(bytes, make([]byte, i)...)

			pools.Release(i, bytes)
		}
	}
}
