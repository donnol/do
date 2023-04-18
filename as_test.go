package do

import "testing"

func TestAs(t *testing.T) {
	var v any = 1
	r := As[int](v)
	if r != 1 {
		t.Errorf("bad case: %d != %d", r, 1)
	}
}
