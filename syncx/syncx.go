package syncx

import (
	"iter"
	"sync"
)

type Map[K comparable, V any] struct {
	sync.Map
}

// old must be comparable
func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.Map.CompareAndDelete(key, old)
}

// old must be comparable
func (m *Map[K, V]) CompareAndSwap(key K, old V, new V) bool {
	return m.Map.CompareAndSwap(key, old, new)
}
func (m *Map[K, V]) Delete(key K) {
	m.Map.Delete(key)
}
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	var tmp any
	if tmp, ok = m.Map.Load(key); ok {
		value = tmp.(V)
	}
	return value, ok
}
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var tmp any
	if tmp, loaded = m.Map.LoadAndDelete(key); loaded {
		value = tmp.(V)
	}
	return value, loaded
}
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	var tmp any
	if tmp, loaded = m.Map.LoadOrStore(key, value); loaded {
		actual = tmp.(V)
	}
	return actual, loaded
}
func (m *Map[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	var tmp any
	if tmp, loaded = m.Map.Swap(key, value); loaded {
		previous = tmp.(V)
	}
	return previous, loaded
}
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Map.Range(func(key, value any) bool {
			return yield(key.(K), value.(V))
		})
	}
}

func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		m.Map.Range(func(key, value any) bool {
			return yield(key.(K))
		})
	}
}
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.Map.Range(func(key, value any) bool {
			return yield(value.(V))
		})
	}
}
