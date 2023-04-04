package do

import "sync"

type Map[K comparable, T any] struct {
	mu sync.Mutex
	m  map[K]T
}

func NewMap[K comparable, T any](size ...int) *Map[K, T] {
	s := 0
	if len(size) > 0 {
		s = size[0]
	}
	return &Map[K, T]{
		m: make(map[K]T, s),
	}
}

func (m *Map[K, T]) Insert(key K, value T) {
	m.mu.Lock()
	m.m[key] = value
	m.mu.Unlock()
}

func (m *Map[K, T]) Lookup(key K) (value T, ok bool) {
	m.mu.Lock()
	value, ok = m.m[key]
	m.mu.Unlock()
	return
}

func (m *Map[K, T]) Get(key K) (value T) {
	value, _ = m.Lookup(key)
	return
}

func (m *Map[K, T]) Range(f func(key K, value T)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range m.m {
		f(k, v)
	}
}

func (m *Map[K, T]) Remove(key K) {
	m.mu.Lock()
	delete(m.m, key)
	m.mu.Unlock()
}
