package dsync

type IntSync struct {
	v  int
	mu RWMutex
}

func NewIntSync(n int) IntSync {
	return IntSync{
		v: n,
	}
}

func (c *IntSync) Get() int {
	defer c.mu.RLock().RUnlock()

	return c.v
}

func (c *IntSync) SetIf(ok func(c int) bool, n int) {
	defer c.mu.Lock().Unlock()

	if !ok(c.v) {
		return
	}

	c.v = n
}
