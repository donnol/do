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

func LogicWPR(f func(C) E) LogicWithoutPR {
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

func (l Logic[P, R]) ToLogic() Logic[P, R] {
	return l
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

type ToLogic[P, R any] interface {
	ToLogic() Logic[P, R]
}

var (
	_ ToLogic[int, string]        = Logic[int, string](nil)
	_ ToLogic[struct{}, string]   = LogicWithoutParam[string](nil)
	_ ToLogic[int, struct{}]      = LogicWithoutResult[int](nil)
	_ ToLogic[struct{}, struct{}] = LogicWithoutPR(nil)
	_ ToLogic[int, string]        = LogicWithoutError[int, string](nil)
	_ ToLogic[struct{}, string]   = LogicWithoutParamError[string](nil)
	_ ToLogic[int, struct{}]      = LogicWithoutResultError[int](nil)
	_ ToLogic[struct{}, struct{}] = LogicWithoutParamResultError(nil)
)

type LogicSet[P, R any] interface {
	func(C, P) (R, E) |
		Logic[P, R] |
		LogicWithoutParam[R] |
		LogicWithoutResult[P] |
		LogicWithoutPR |
		LogicWithoutError[P, R] |
		LogicWithoutParamError[R] |
		LogicWithoutResultError[P] |
		LogicWithoutParamResultError
}
