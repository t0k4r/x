package cachex

import (
	"iter"
	"time"

	"github.com/t0k4r/x/iterx"
	"github.com/t0k4r/x/syncx"
)

type Cache[K comparable, V any] struct {
	mp syncx.Map[K, item[V]]
}

type item[T any] struct {
	value T
	timer *time.Timer
}

func (c *Cache[K, V]) Delete(key K) (value V, loaded bool) {
	var it item[V]
	if it, loaded = c.mp.LoadAndDelete(key); loaded {
		it.timer.Stop()
		value = it.value
	}
	return value, loaded
}

func (c *Cache[K, V]) Load(key K) (value V, ok bool) {
	var it item[V]
	if it, ok = c.mp.Load(key); ok {
		value = it.value
	}
	return value, ok
}

func (c *Cache[K, V]) Store(key K, value V, duration time.Duration) (previous V, loaded bool) {
	var oldIt item[V]
	newIt := item[V]{value, time.AfterFunc(duration, func() { c.mp.Delete(key) })}
	if oldIt, loaded = c.mp.LoadOrStore(key, newIt); loaded {
		oldIt.timer.Stop()
		previous = oldIt.value
		c.mp.Swap(key, newIt)
	}
	return previous, loaded
}
func (c *Cache[K, V]) All() iter.Seq2[K, V] {
	return iterx.Map2(c.mp.All(), func(k K, it item[V]) (K, V) { return k, it.value })
}

func (c *Cache[K, V]) Keys() iter.Seq[K] {
	return c.mp.Keys()
}
func (c *Cache[K, V]) Values() iter.Seq[V] {
	return iterx.Map(c.mp.Values(), func(it item[V]) V { return it.value })
}
