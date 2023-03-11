package dsync_test

import (
	"fmt"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
TestPoolBuffer - this test may fail, because sync.Pool may decide not to use pre-allocated buffer. In this case test may
be run several times and one of tries should end successfully.
*/
func TestPoolBuffer(t *testing.T) {
	var inputs = []struct {
		SizeGet int
	}{
		{
			SizeGet: dsync.PoolBufferSize64B,
		},
		{
			SizeGet: dsync.PoolBufferSize128B,
		},
		{
			SizeGet: dsync.PoolBufferSize256B,
		},
		{
			SizeGet: dsync.PoolBufferSize512B,
		},
		{
			SizeGet: dsync.PoolBufferSizeMax,
		},
	}

	var buffersPointersBySize = make(map[int]string)

	for _, input := range inputs {
		var input = input

		t.Run(fmt.Sprintf("size_%dB", input.SizeGet), func(t *testing.T) {
			// Inflate buffers pool.
			var bufferNew = dsync.PoolBufferGet(input.SizeGet)
			require.EqualValues(t, input.SizeGet, bufferNew.Cap())

			func() {
				defer dsync.PoolBufferPut(bufferNew)
				buffersPointersBySize[bufferNew.Cap()] = fmt.Sprintf("%p", bufferNew)
			}()

			// Check if buffers pointers are the same.
			var bufferOld = dsync.PoolBufferGet(input.SizeGet)
			require.EqualValues(t, input.SizeGet, bufferOld.Cap())

			func() {
				defer dsync.PoolBufferPut(bufferOld)
				require.EqualValues(t, buffersPointersBySize[bufferOld.Cap()], fmt.Sprintf("%p", bufferOld))
			}()
		})
	}
}
