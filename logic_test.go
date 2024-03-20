package do

import (
	"context"
	"testing"
)

func TestToLogic(t *testing.T) {
	var text = "test"

	logic := func(ctx context.Context, p int) (string, error) { return text, nil }

	r := Must1(LogicFrom(logic).ToLogic()(context.Background(), 1))
	Assert(t, r, text)

	r1 := Must1(logicHelper[int, string](Logic[int, string](logic))(context.Background(), 2))
	Assert(t, r1, text)
}

func logicHelper[P, R any](logic ToLogic[P, R]) Logic[P, R] {
	return logic.ToLogic()
}
