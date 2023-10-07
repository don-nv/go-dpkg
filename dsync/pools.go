package dsync

import (
	"sync"
)

type BytesPools struct {
	byID map[int]BytesPool
	mu   RWMutex
}

func NewPoolsBytes() BytesPools {
	return BytesPools{
		byID: make(map[int]BytesPool),
		mu:   NewRWMutex(),
	}
}

func (p *BytesPools) Acquire(id int) []byte {
	var (
		pool BytesPool
		is   bool
	)

	p.mu.RLockF(func() { pool, is = p.byID[id] })

	if !is {
		p.mu.LockF(func() {
			pool, is = p.byID[id]
			if !is {
				pool = NewPoolBytes()
				p.byID[id] = pool
			}
		})
	}

	return pool.Acquire()
}

func (p *BytesPools) Release(id int, bytes []byte) {
	p.mu.RLockF(func() {
		pool, ok := p.byID[id]
		if ok {
			pool.Release(bytes)
		}
	})
}

type BytesPool struct {
	pool *sync.Pool
	cap  *IntSync
}

func NewPoolBytes() BytesPool {
	var capSync IntSync

	return BytesPool{
		pool: &sync.Pool{
			New: func() interface{} { return make([]byte, 0, capSync.Get()) },
		},
		cap: &capSync,
	}
}

func (p BytesPool) Acquire() []byte { return p.pool.Get().([]byte) }

func (p BytesPool) Release(bytes []byte) {
	if n := cap(bytes); n > p.cap.Get() {
		p.cap.SetIf(func(c int) bool { return n > c }, n)
	}

	bytes = bytes[:0]

	p.pool.Put(bytes)
}
