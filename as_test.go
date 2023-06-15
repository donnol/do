package do

import (
	"fmt"
	"testing"
)

func TestAs(t *testing.T) {
	var v any = 1
	r := As[int](v)
	if r != 1 {
		t.Errorf("bad case: %d != %d", r, 1)
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				if fmt.Sprintf("%v", r) != "interface conversion: interface {} is int, not string" {
					t.Errorf("bad case: %v", r)
				}
			}
		}()

		As[string](v)
	}()
}

func TestAsOk(t *testing.T) {
	var v any = 1
	r, ok := AsOk[int](v)
	if !ok || r != 1 {
		t.Errorf("bad case: %d != %d", r, 1)
	}

	{
		_, ok := AsOk[string](v)
		if ok {
			t.Errorf("bad case, v is not a string")
		}
	}
}
