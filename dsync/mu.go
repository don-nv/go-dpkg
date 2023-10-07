package dsync

import "sync"

type RWMutex struct{ mu sync.RWMutex }

func NewRWMutex() RWMutex { return RWMutex{} }

func (m *RWMutex) LockF(f func()) *RWMutex  { defer m.Lock().Unlock(); f(); return m }
func (m *RWMutex) RLockF(f func()) *RWMutex { defer m.RLock().RUnlock(); f(); return m }

func (m *RWMutex) Lock() *RWMutex    { m.mu.Lock(); return m }
func (m *RWMutex) Unlock() *RWMutex  { m.mu.Unlock(); return m }
func (m *RWMutex) RLock() *RWMutex   { m.mu.RLock(); return m }
func (m *RWMutex) RUnlock() *RWMutex { m.mu.RUnlock(); return m }
