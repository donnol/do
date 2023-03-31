package do

import (
	"strings"
	"testing"
)

func TestFuncName(t *testing.T) {
	r := FuncName(1, true)
	want := "do/func_name_test.go:9 github.com/donnol/do.TestFuncName"
	if !strings.Contains(r, want) {
		t.Errorf("bad case, %s didn't contains %s", r, want)
	}

	{
		r := FuncName(1, false)
		want := "github.com/donnol/do.TestFuncName"
		if r != want {
			t.Errorf("bad case, %s != %s", r, want)
		}
	}
}
