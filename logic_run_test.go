package do

import (
	"context"
	"testing"
)

func TestRunIf(t *testing.T) {
	r := Must1(RunIf[int, string](1 != 0, context.Background(), 1, logic))
	Assert(t, r, "")
}
