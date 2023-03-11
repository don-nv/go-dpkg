package dcache

import (
	"container/list"
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	options[K, V]

	Set(key K, val V, ttl time.Duration)
	Get(key K) (V, bool)
	DeleteExpiredCache()
	Len() int
	InvalidateKey(key K)

	removeItem(e *list.Element)
	removeBack()
	removeBackIfExpired()
}

const weekDuration = time.Hour * 24 * 7

func NewCache[K comparable, V any]() Cache[K, V] {
	return &cache[K, V]{
		data:     make(map[K]*list.Element),
		items:    list.New(),
		capacity: -1,
		ttl:      weekDuration,
	}
}

// Set - ttl == 0 equals time.Hour * 24 * 7 (one week).
func (c *cache[K, V]) Set(key K, val V, ttl time.Duration) {
	c.mu.Lock()

	if ttl <= 0 {
		ttl = c.ttl
	}

	now := time.Now()
	if el, ok := c.data[key]; ok {
		el.Value.(*item[K, V]).val = val
		el.Value.(*item[K, V]).expiredAt = now.Add(ttl)

		c.items.MoveToFront(el)

		c.mu.Unlock()
		return
	}

	el := &item[K, V]{
		key:       key,
		val:       val,
		expiredAt: now.Add(ttl),
	}

	element := c.items.PushFront(el)
	c.data[key] = element

	c.removeBackIfExpired()

	if c.capacity >= 0 && c.items.Len() > c.capacity {
		c.removeBack()
	}

	c.mu.Unlock()
}

func (c *cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()

	result := *new(V)
	if el, ok := c.data[key]; ok {
		if time.Now().After(el.Value.(*item[K, V]).expiredAt) {
			c.mu.RUnlock()
			return result, false
		}

		c.items.MoveToFront(el)

		c.mu.RUnlock()
		return el.Value.(*item[K, V]).val, true
	}

	c.mu.RUnlock()
	return result, false
}

func (c *cache[K, V]) DeleteExpiredCache() {
	c.mu.Lock()

	for it := c.items.Front(); it != nil; it = it.Next() {
		if time.Now().After(it.Value.(*item[K, V]).expiredAt) {
			c.removeItem(it)
		}
	}

	c.mu.Unlock()
}

func (c *cache[K, V]) Len() int {
	c.mu.RLock()
	lens := len(c.data)
	c.mu.RUnlock()
	return lens
}

func (c *cache[K, V]) InvalidateKey(key K) {
	c.mu.Lock()
	if el, ok := c.data[key]; ok {
		c.removeItem(el)
	}
	c.mu.Unlock()
}

func (c *cache[K, V]) removeItem(e *list.Element) {
	c.items.Remove(e)
	kv := e.Value.(*item[K, V]) //nolint:errcheck
	delete(c.data, kv.key)
}

func (c *cache[K, V]) removeBack() {
	element := c.items.Back()
	if element != nil {
		c.removeItem(element)
	}
}

func (c *cache[K, V]) removeBackIfExpired() {
	element := c.items.Back()
	if element != nil && time.Now().After(element.Value.(*item[K, V]).expiredAt) {
		c.removeItem(element)
	}
}

type cache[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]*list.Element

	capacity int
	ttl      time.Duration
	items    *list.List
}

type item[K comparable, V any] struct {
	expiredAt time.Time
	key       K
	val       V
}

type options[K comparable, V any] interface {
	WithCapacity(capacity int) Cache[K, V]
	WithDuration(ttl time.Duration) Cache[K, V]
}

func (c *cache[K, V]) WithCapacity(capacity int) Cache[K, V] {
	c.capacity = capacity
	return c
}

func (c *cache[K, V]) WithDuration(ttl time.Duration) Cache[K, V] {
	c.ttl = ttl
	return c
}
