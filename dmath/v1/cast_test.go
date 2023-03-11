package dmath_test

import (
	"github.com/don-nv/go-dpkg/dmath/v1"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestCast(t *testing.T) {
	_, ok := dmath.Cast[int, int32](math.MaxInt)
	require.False(t, ok)

	_, ok = dmath.Cast[int, int32](math.MinInt)
	require.False(t, ok)

	_, ok = dmath.Cast[int32, int](math.MaxInt32)
	require.True(t, ok)

	_, ok = dmath.Cast[int32, int](math.MinInt32)
	require.True(t, ok)

	_, ok = dmath.Cast[int, int64](math.MaxInt64)
	require.True(t, ok)

	_, ok = dmath.Cast[int, int64](math.MinInt64)
	require.True(t, ok)

	_, ok = dmath.Cast[int64, int](math.MaxInt64)
	require.True(t, ok)

	_, ok = dmath.Cast[int64, int](math.MinInt64)
	require.True(t, ok)

	_, ok = dmath.Cast[int32, int64](math.MaxInt32)
	require.True(t, ok)

	_, ok = dmath.Cast[int32, int64](math.MinInt32)
	require.True(t, ok)

	_, ok = dmath.Cast[int64, int32](math.MaxInt64)
	require.False(t, ok)

	_, ok = dmath.Cast[int64, int32](math.MinInt64)
	require.False(t, ok)
}
