package xmap

import (
	"sync"
)

type XMap[K comparable, V any] struct {
	m *sync.Map
}

func NewXMap[K comparable, V any]() *XMap[K, V] {
	return &XMap[K, V]{
		m: &sync.Map{},
	}
}

func (m *XMap[K, V]) Store(key K, value V) {
	if m == nil || m.m == nil {
		return
	}

	m.m.Store(key, value)
}

func (m *XMap[K, V]) Load(key K) (V, bool) {
	if m == nil || m.m == nil {
		var v V
		return v, false
	}

	val, ok := m.m.Load(key)
	if !ok {
		var temp V
		return temp, ok
	}
	return val.(V), ok
}

func (m *XMap[K, V]) Has(key K) bool {
	if m == nil || m.m == nil {
		return false
	}

	_, ok := m.m.Load(key)
	return ok
}

func (m *XMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	if m == nil || m.m == nil {
		var v V
		return v, false
	}

	val, ok := m.m.LoadOrStore(key, value)
	if !ok {
		return value, ok
	}
	return val.(V), ok
}

func (m *XMap[K, V]) LoadAndStore(key K, value V) (V, bool) {
	if m == nil || m.m == nil {
		var v V
		return v, false
	}

	val, ok := m.m.Load(key)
	if ok {
		m.m.Swap(key, value)
		return val.(V), ok
	}
	m.m.Store(key, value)
	return value, false
}

func (m *XMap[K, V]) LoadAndDelete(key K) (V, bool) {
	if m == nil || m.m == nil {
		var v V
		return v, false
	}

	val, ok := m.m.LoadAndDelete(key)
	if !ok {
		var temp V
		return temp, ok
	}
	return val.(V), ok
}

func (m *XMap[K, V]) Delete(key K) {
	if m == nil || m.m == nil {
		return
	}

	m.m.Delete(key)
}

func (m *XMap[K, V]) Swap(key K, value V) (V, bool) {
	if m == nil || m.m == nil {
		var v V
		return v, false
	}

	val, ok := m.m.Swap(key, value)
	if !ok {
		var v V
		return v, ok
	}
	return val.(V), ok
}

func (m *XMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	if m == nil || m.m == nil {
		return false
	}

	return m.m.CompareAndSwap(key, old, new)
}

func (m *XMap[K, V]) Range(fn func(key K, value V) bool) {
	if m == nil || m.m == nil || fn == nil {
		return
	}

loopLabel:
	for key, value := range m.m.Range {
		k, ok1 := key.(K)
		v, ok2 := value.(V)
		switch {
		case ok1 && ok2:
			if !fn(k, v) {
				break loopLabel
			}
		case ok1:
			var temp V
			if !fn(k, temp) {
				break loopLabel
			}
		case ok2:
			var temp K
			if !fn(temp, v) {
				break loopLabel
			}
		}
	}
}

func (m *XMap[K, V]) RangeAndDelete(fn func(key K, value V) bool) {
	if m == nil || m.m == nil || fn == nil {
		return
	}

	m.Range(func(key K, value V) bool {
		if fn(key, value) {
			m.m.Delete(key)
		}
		return true
	})
}

func (m *XMap[K, V]) DeleteKeys(keys []K) {
	if m == nil || m.m == nil || len(keys) == 0 {
		return
	}

	for _, key := range keys {
		m.m.Delete(key)
	}
}

func (m *XMap[K, V]) DeleteWhere(fn func(key K, value V) bool) {
	if m == nil || m.m == nil || fn == nil {
		return
	}

	m.Range(func(key K, value V) bool {
		if fn(key, value) {
			m.m.Delete(key)
		}
		return true
	})
}

func (m *XMap[_, _]) Size() int {
	if m == nil || m.m == nil {
		return 0
	}

	var count int
	for range m.m.Range {
		count++
	}
	return count
}

func (m *XMap[K, V]) Clear() {
	if m == nil || m.m == nil {
		return
	}

	m.m.Clear()
}
