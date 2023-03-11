package dsync

import "sync"

type RWMutex struct{ mu sync.RWMutex }

func (m *RWMutex) LockF(f func()) *RWMutex  { defer m.Lock().Unlock(); f(); return m }
func (m *RWMutex) RLockF(f func()) *RWMutex { defer m.RLock().RUnlock(); f(); return m }

func (m *RWMutex) Lock() Unlocker   { m.mu.Lock(); return Unlocker{m} }
func (m *RWMutex) RLock() RUnlocker { m.mu.RLock(); return RUnlocker{m} }

// Unlocker - must not be initialized directly, but with functions or method returning instance.
type Unlocker struct {
	mu *RWMutex
}

func (u Unlocker) Unlock() {
	u.mu.mu.Unlock()
}

// RUnlocker - must not be initialized directly, but with functions or method returning instance.
type RUnlocker struct {
	mu *RWMutex
}

func (u RUnlocker) RUnlock() {
	u.mu.mu.RUnlock()
}
