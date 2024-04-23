package do

func RunIf[P, R any, F LogicSet[P, R]](cond bool, ctx C, param P, f F) (r R, err E) {
	if !cond {
		return
	}

	var l Logic[P, R]
	switch lc := any(f).(type) {
	case func(C, P) (R, E):
		l = lc
	case ToLogic[P, R]:
		l = lc.ToLogic()
	default:
		panic("unsupport func")
	}

	return l(ctx, param)
}
