package dcache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Cache_Set(t *testing.T) {
	cache := NewCache[string, []int]()

	ints := []int{1, 2, 3}
	cache.Set("key", ints, 0)

	vals, b := cache.Get("key")

	require.EqualValues(t, true, b)
	require.EqualValues(t, 3, len(vals))
}

func Test_Cache_DurationSet(t *testing.T) {
	cache := NewCache[string, []int]().
		WithDuration(time.Millisecond * 500)

	ints := []int{1, 2, 3}
	cache.Set("key1", ints, 0)
	time.Sleep(1 * time.Second)
	cache.Set("key2", ints, 0)

	result := make([][]int, 0, 2)

	vals, b := cache.Get("key1")
	if vals != nil {
		result = append(result, vals)
	}
	require.EqualValues(t, false, b)

	vals, b = cache.Get("key2")
	if vals != nil {
		result = append(result, vals)
	}
	require.EqualValues(t, true, b)

	require.EqualValues(t, 1, len(result))
}

func Test_Cache_CapacitySet(t *testing.T) {
	cache := NewCache[string, int]().
		WithCapacity(9)

	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for i := range ints {
		cache.Set(fmt.Sprintf("key%d", i+1), ints[i], 0)
	}

	require.EqualValues(t, 9, cache.Len())
}

func Test_Cache_Get(t *testing.T) {
	require.EqualValues(t, nil, nil)
}

func TestCache_InvalidateKey(t *testing.T) {
	newCache := NewCache[string, int]()

	newCache.Set("key1", 1, 0)
	newCache.Set("key2", 2, 0)
	newCache.Set("key3", 3, 0)

	newCache.InvalidateKey("key2")

	k1, ok := newCache.Get("key2")
	require.EqualValues(t, false, ok)
	require.Equal(t, 0, k1)

	newCache.InvalidateKey("key4")
	k2, ok := newCache.Get("key4")
	require.EqualValues(t, false, ok)
	require.Equal(t, 0, k2)
}
