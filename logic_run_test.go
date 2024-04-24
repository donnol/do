package do

import (
	"context"
	"testing"
)

func TestRunIf(t *testing.T) {
	{
		r := Must1(RunIf[int, string](1 != 0, context.Background(), 1, logic))
		Assert(t, r, "")
	}
	{
		r := Must1(RunIf[struct{}, string](1 != 0, context.Background(), struct{}{}, LogicWP(logicWP)))
		Assert(t, r, "")
	}
	{
		r := Must1(RunIf[int, struct{}](1 != 0, context.Background(), 1, LogicWR(logicWR)))
		Assert(t, r, struct{}{})
	}
	{
		r := Must1(RunIf[struct{}, struct{}](1 != 0, context.Background(), struct{}{}, LogicWPR(logicWPR)))
		Assert(t, r, struct{}{})
	}
	{
		r := Must1(RunIf[int, string](1 != 0, context.Background(), 1, LogicWE(logicWE)))
		Assert(t, r, "")
	}
	{
		r := Must1(RunIf[struct{}, string](1 != 0, context.Background(), struct{}{}, LogicWPE(logicWPE)))
		Assert(t, r, "")
	}
	{
		r := Must1(RunIf[int, struct{}](1 != 0, context.Background(), 1, LogicWRE(logicWRE)))
		Assert(t, r, struct{}{})
	}
	{
		r := Must1(RunIf[struct{}, struct{}](1 != 0, context.Background(), struct{}{}, LogicWPRE(logicWPRE)))
		Assert(t, r, struct{}{})
	}
}

func TestRunLogicIf(t *testing.T) {
	{
		r := Must1(RunLogicIf(1 != 0, context.Background(), 1, LogicFrom(logic)))
		Assert(t, r, "")
	}
	{
		r := Must1(RunLogicIf(1 != 0, context.Background(), struct{}{}, LogicFromWP(logicWP)))
		Assert(t, r, "")
	}
	{
		r := Must1(RunLogicIf(1 != 0, context.Background(), 1, LogicFromWR(logicWR)))
		Assert(t, r, struct{}{})
	}
	{
		r := Must1(RunLogicIf(1 != 0, context.Background(), struct{}{}, LogicFromWPR(logicWPR)))
		Assert(t, r, struct{}{})
	}
	{
		r := Must1(RunLogicIf(1 != 0, context.Background(), 1, LogicFromWE(logicWE)))
		Assert(t, r, "")
	}
	{
		r := Must1(RunLogicIf(1 != 0, context.Background(), struct{}{}, LogicFromWPE(logicWPE)))
		Assert(t, r, "")
	}
	{
		r := Must1(RunLogicIf(1 != 0, context.Background(), 1, LogicFromWRE(logicWRE)))
		Assert(t, r, struct{}{})
	}
	{
		r := Must1(RunLogicIf(1 != 0, context.Background(), struct{}{}, LogicFromWPRE(logicWPRE)))
		Assert(t, r, struct{}{})
	}
}
