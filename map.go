package do

import (
	"log"
	"sync"
	"time"
)

var (
	watchInterval = 15 * time.Second
)

type Map[K comparable, T any] struct {
	mu       sync.Mutex
	m        map[K]*mapEntry[T]
	expireds map[K]*mapEntry[T]
}

type mapEntry[T any] struct {
	value    T
	deadline time.Time
}

func NewMap[K comparable, T any](size ...int) *Map[K, T] {
	s := 0
	if len(size) > 0 {
		s = size[0]
	}
	m := &Map[K, T]{
		m:        make(map[K]*mapEntry[T], s),
		expireds: make(map[K]*mapEntry[T], s),
	}

	go m.watch(watchInterval)

	return m
}

type KeyOption struct {
	Timeout time.Duration // key will expire when timeout
}

func (m *Map[K, T]) Insert(key K, value T, opts ...*KeyOption) {
	var deadline time.Time
	if len(opts) > 0 {
		if opts[0].Timeout != 0 {
			deadline = time.Now().Add(opts[0].Timeout)
		}
	}

	e := &mapEntry[T]{
		value:    value,
		deadline: deadline,
	}

	m.mu.Lock()
	m.m[key] = e
	if !deadline.IsZero() {
		m.expireds[key] = e
	}
	m.mu.Unlock()
}

func (m *Map[K, T]) Lookup(key K) (value T, ok bool) {
	m.mu.Lock()
	v, ok := m.m[key]
	m.mu.Unlock()

	// 检查是否过期
	if !ok || IsExpired(v.deadline, time.Now()) {
		return value, false
	}

	return v.value, true
}

func (m *Map[K, T]) Get(key K) (value T) {
	value, _ = m.Lookup(key)
	return
}

func (m *Map[K, T]) Range(f func(key K, value T)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range m.m {
		if IsExpired(v.deadline, time.Now()) {
			continue
		}
		f(k, v.value)
	}
}

func (m *Map[K, T]) Remove(key K) {
	m.mu.Lock()
	delete(m.m, key)
	m.mu.Unlock()
}

func (m *Map[K, T]) watch(watchInterval time.Duration) {
	// 每隔一段时间检查key是否过期
	for {
		time.Sleep(watchInterval)

		func() {
			m.mu.Lock()
			defer m.mu.Unlock()

			now := time.Now()
			for k, v := range m.expireds {
				if IsExpired(v.deadline, now) {
					delete(m.m, k)
					delete(m.expireds, k)

					log.Printf("delete key %v\n at %v", k, now)
				}
			}
		}()
	}
}

// ValueAs get value by key from *Map[K, any], and assert value type to T
func ValueAs[K comparable, T any](m *Map[K, any], key K) T {
	return As[T](m.Get(key))
}
