package syncx

import (
	"sync"
)

type Map[K, V any] struct {
	sync.Map
}

func NewMap[K, V any]() *Map[K, V] {
	return &Map[K, V]{
		Map: sync.Map{},
	}
}

func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.Map.CompareAndDelete(key, old)
}
func (m *Map[K, V]) CompareAndSwap(key K, old V, new V) bool {
	return m.Map.CompareAndSwap(key, old, new)
}
func (m *Map[K, V]) Delete(key K) {
	m.Map.Delete(key)
}
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	var tmp any
	tmp, ok = m.Map.Load(key)
	if !ok {
		value = tmp.(V)
	}
	return
}
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var tmp any
	tmp, loaded = m.Map.LoadAndDelete(key)
	if !loaded {
		value = tmp.(V)
	}
	return
}
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.Map.LoadOrStore(key, value)
	return a.(V), loaded
}
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.Map.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}
func (m *Map[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	var tmp any
	tmp, loaded = m.Map.Swap(key, value)
	if loaded {
		previous = tmp.(V)
	}
	return
}
