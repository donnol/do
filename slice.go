package do

import (
	"sync"
)

type Slice[T any] struct {
	mu    sync.Mutex
	items []T
}

func NewSlice[T any](lenAndCap ...int) *Slice[T] {
	l, cap := getLenAndCap(lenAndCap...)

	return &Slice[T]{
		items: make([]T, l, cap),
	}
}

func (s *Slice[T]) Index(i int) T {
	s.mu.Lock()
	item := s.items[i]
	s.mu.Unlock()

	return item
}

func (s *Slice[T]) Append(values ...T) {
	s.mu.Lock()
	s.items = append(s.items, values...)
	s.mu.Unlock()
}

func (s *Slice[T]) Range(f func(item T, index int)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, v := range s.items {
		f(v, i)
	}
}

func (s *Slice[T]) Reset(lenAndCap ...int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, c := getLenAndCap(lenAndCap...)
	s.items = make([]T, l, c)
}

func getLenAndCap(lenAndCap ...int) (l, c int) {
	if len(lenAndCap) >= 1 {
		l = lenAndCap[0]
	}
	if len(lenAndCap) >= 2 {
		c = lenAndCap[1]
	}
	return
}
