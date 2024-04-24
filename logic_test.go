package do

import (
	"context"
	"net/http"
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

func TestRouteFromLogic(t *testing.T) {
	route := routeFromLogic[int, string](http.MethodPost, "/abc", "abc api", logic)
	Assert(t, route.Method, http.MethodPost)
	Assert(t, route.Path, "/abc")

	{
		route := routeFromLogic[struct{}, string](http.MethodPost, "/abc", "abc api", LogicWP(logicWP))
		Assert(t, route.Path, "/abc")
	}
}

func logic(ctx context.Context, p int) (r string, err error) { return }
func logicWP(ctx context.Context) (r string, err error)      { return }
func logicWR(ctx context.Context, p int) (err error)         { return }
func logicWPR(ctx context.Context) (err error)               { return }
func logicWE(ctx context.Context, p int) (r string)          { return }
func logicWPE(ctx context.Context) (r string)                { return }
func logicWRE(ctx context.Context, p int)                    { return }
func logicWPRE(ctx context.Context)                          { return }

func logicHelper[P, R any](logic ToLogic[P, R]) Logic[P, R] {
	return logic.ToLogic()
}

// 在1.20版本，将L顺序提前到P和R之前也会出现类型推导错误
func routeFromLogic[P, R any, L LogicSet[P, R]](method, path, comment string, lf L) *Route[struct{}] {
	var l Logic[P, R]
	switch lc := any(lf).(type) {
	case func(C, P) (R, E):
		l = lc
	case ToLogic[P, R]:
		l = lc.ToLogic()
	default:
		panic("unsupport func")
	}
	_ = l

	return &Route[struct{}]{
		Method:  method,
		Path:    path,
		Comment: comment,
	}
}
