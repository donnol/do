package do

import "testing"

func TestSlice(t *testing.T) {
	s := NewSlice[int](0, 10)
	s.Append(1)

	r := s.Index(0)
	if r != 1 {
		t.Errorf("bad case: %v != %v", r, 1)
	}

	s.Append(2)
	r = s.Index(1)
	if r != 2 {
		t.Errorf("bad case: %v != %v", r, 2)
	}

	s.Range(func(item, index int) {
		if index == 0 && item != 1 {
			t.Errorf("bad case: %v != %v", item, 1)
		}
		if index == 1 && item != 2 {
			t.Errorf("bad case: %v != %v", item, 2)
		}
	})
}
