package do

var (
	EmptyLogicFunc = func(C, struct{}) (r struct{}, err E) { return }
	NilLogicFunc   Logic[struct{}, struct{}]
)

type (
	Logic[P, R any]                func(C, P) (R, E)
	LogicWithoutParam[R any]       func(C) (R, E)
	LogicWithoutResult[P any]      func(C, P) E
	LogicWithoutPR                 func(C) E
	LogicWithoutError[P, R any]    func(C, P) R
	LogicWithoutParamError[R any]  func(C) R
	LogicWithoutResultError[P any] func(C, P)
	LogicWithoutParamResultError   func(C)
)

func LogicFrom[P, R any](f func(C, P) (R, E)) Logic[P, R] {
	return f
}

func LogicWP[R any](f func(C) (R, E)) LogicWithoutParam[R] {
	return f
}

func LogicWR[P any](f func(C, P) E) LogicWithoutResult[P] {
	return f
}

func LogicWPR[P any](f func(C) E) LogicWithoutPR {
	return f
}

func LogicWE[P, R any](f func(C, P) R) LogicWithoutError[P, R] {
	return f
}

func LogicWPE[R any](f func(C) R) LogicWithoutParamError[R] {
	return f
}

func LogicWRE[P any](f func(C, P)) LogicWithoutResultError[P] {
	return f
}

func LogicWPRE(f func(C)) LogicWithoutParamResultError {
	return f
}

func (l LogicWithoutParam[R]) ToLogic() Logic[struct{}, R] {
	return func(ctx C, p struct{}) (r R, err E) {
		return l(ctx)
	}
}

func (l LogicWithoutResult[P]) ToLogic() Logic[P, struct{}] {
	return func(ctx C, p P) (r struct{}, err E) {
		err = l(ctx, p)
		return
	}
}

func (l LogicWithoutPR) ToLogic() Logic[struct{}, struct{}] {
	return func(ctx C, p struct{}) (r struct{}, err E) {
		err = l(ctx)
		return
	}
}

func (l LogicWithoutError[P, R]) ToLogic() Logic[P, R] {
	return func(ctx C, p P) (r R, err E) {
		r = l(ctx, p)
		return
	}
}

func (l LogicWithoutParamError[R]) ToLogic() Logic[struct{}, R] {
	return func(ctx C, p struct{}) (r R, err E) {
		r = l(ctx)
		return
	}
}

func (l LogicWithoutResultError[P]) ToLogic() Logic[P, struct{}] {
	return func(ctx C, p P) (r struct{}, err E) {
		l(ctx, p)
		return
	}
}

func (l LogicWithoutParamResultError) ToLogic() Logic[struct{}, struct{}] {
	return func(ctx C, p struct{}) (r struct{}, err E) {
		l(ctx)
		return
	}
}
