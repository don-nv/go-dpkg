package dslice_test

import (
	"github.com/don-nv/go-dpkg/dslice/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSlice_Append(t *testing.T) {
	var s = dslice.NewSlice[int]()

	for i := 0; i < 1024; i++ {
		s.Append(i)
	}

	s.Range(func(i int, n int) {
		require.EqualValues(t, i, n)
	})
}
