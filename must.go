package do

import "fmt"

// Must args的最后一个是非空错误时panic，无结果返回
func Must(args ...any) {
	mustCheckArgs(1, args...)
}

// Must1 args的最后一个是非空错误时panic，返回1个结果
func Must1[T any](args ...any) (r T) {
	l := len(args)
	mustCheckArgs(2, args...)

	r = args[l-2].(T)

	return
}

// Must2 args的最后一个是非空错误时panic，返回2个结果
func Must2[T1, T2 any](args ...any) (r1 T1, r2 T2) {
	l := len(args)
	mustCheckArgs(3, args...)

	r2 = args[l-2].(T2)
	r1 = args[l-3].(T1)

	return
}

// Must3 args的最后一个是非空错误时panic，返回3个结果
func Must3[T1, T2, T3 any](args ...any) (r1 T1, r2 T2, r3 T3) {
	l := len(args)
	mustCheckArgs(4, args...)

	r3 = args[l-2].(T3)
	r2 = args[l-3].(T2)
	r1 = args[l-4].(T1)

	return
}

// Must4 args的最后一个是非空错误时panic，返回4个结果
func Must4[T1, T2, T3, T4 any](args ...any) (r1 T1, r2 T2, r3 T3, r4 T4) {
	l := len(args)
	mustCheckArgs(5, args...)

	r4 = args[l-2].(T4)
	r3 = args[l-3].(T3)
	r2 = args[l-4].(T2)
	r1 = args[l-5].(T1)

	return
}

// Must5 args的最后一个是非空错误时panic，返回5个结果
func Must5[T1, T2, T3, T4, T5 any](args ...any) (r1 T1, r2 T2, r3 T3, r4 T4, r5 T5) {
	l := len(args)
	mustCheckArgs(6, args...)

	r5 = args[l-2].(T5)
	r4 = args[l-3].(T4)
	r3 = args[l-4].(T3)
	r2 = args[l-5].(T2)
	r1 = args[l-6].(T1)

	return
}

func mustCheckArgs(wantLength int, args ...any) {
	if wantLength <= 0 {
		return
	}

	l := len(args)
	if l != wantLength {
		panic(NewError(efrom(fmt.Errorf("args length is not equal %d", wantLength))))
	}

	if v, ok := args[l-1].(error); ok && v != nil {
		panic(NewError(efrom(v)))
	}
}
