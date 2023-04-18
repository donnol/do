package do

import "testing"

type F struct {
	R int

	// If use next line, it will complain with 'T does not satisfy comparable'
	// A any
}

// IsFComparable make sure F is a comparable type in compile time
func IsFComparable[T F]() {
	IsComparable[T]()
}

func TestComparable(t *testing.T) {
	f1 := F{R: 1}
	f2 := F{R: 1}
	if f1 != f2 {
		t.Errorf("bad case: %v != %v", f1, f2)
	}
}
