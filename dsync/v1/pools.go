package dsync

import (
	"bytes"
	"sync"
)

const (
	PoolBufferSize64B  = 64
	PoolBufferSize128B = 128
	PoolBufferSize256B = 256
	PoolBufferSize512B = 512
	PoolBufferSizeMax  = 1024
)

var poolBuffs = func() map[int]*Pool[*bytes.Buffer] {
	var (
		sizes = []int{
			PoolBufferSize64B,
			PoolBufferSize128B,
			PoolBufferSize256B,
			PoolBufferSize512B,
			PoolBufferSizeMax,
		}
		buffs = make(map[int]*Pool[*bytes.Buffer], len(sizes))
	)

	for _, n := range sizes {
		var n = n

		buffs[n] = NewPool(
			func() *bytes.Buffer {
				return bytes.NewBuffer(make([]byte, 0, n))
			},
			OptionPoolWithOnPut(func(t *bytes.Buffer) {
				t.Reset()
			}),
		)
	}

	return buffs
}()

func PoolBufferGet(size int) *bytes.Buffer {
	if size <= PoolBufferSize64B {
		return poolBuffs[PoolBufferSize64B].Get()
	}

	if size <= PoolBufferSize128B {
		return poolBuffs[PoolBufferSize128B].Get()
	}

	if size <= PoolBufferSize256B {
		return poolBuffs[PoolBufferSize256B].Get()
	}

	if size <= PoolBufferSize512B {
		return poolBuffs[PoolBufferSize512B].Get()
	}

	return poolBuffs[PoolBufferSizeMax].Get()
}

func PoolBufferPut(buf *bytes.Buffer) {
	var n = buf.Cap()

	if n <= PoolBufferSize64B {
		poolBuffs[PoolBufferSize64B].Put(buf)

		return
	}

	if n <= PoolBufferSize128B {
		poolBuffs[PoolBufferSize128B].Put(buf)

		return
	}

	if n <= PoolBufferSize256B {
		poolBuffs[PoolBufferSize256B].Put(buf)

		return
	}

	if n <= PoolBufferSize512B {
		poolBuffs[PoolBufferSize512B].Put(buf)

		return
	}

	poolBuffs[PoolBufferSizeMax].Put(buf)
}

type Pool[T any] struct {
	pool  sync.Pool
	onPut func(t T)
}

func NewPool[T any](new func() T, options ...PoolOption[T]) *Pool[T] {
	pool := Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return new()
			},
		},
		onPut: func(T) {},
	}

	for _, option := range options {
		option(&pool)
	}

	return &pool
}

func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(t T) {
	p.onPut(t)
	p.pool.Put(t)
}

type PoolOption[T any] func(p *Pool[T])

func OptionPoolWithOnPut[T any](f func(t T)) PoolOption[T] {
	return func(p *Pool[T]) {
		p.onPut = f
	}
}
