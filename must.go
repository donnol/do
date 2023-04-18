package do

import "fmt"

// Must panic if err is not nill
func Must(err error) {
	mustCheckError(err)
}

// Must1 panic if err is not nill，or return 1 result
func Must1[T any](a1 T, err error) T {
	mustCheckError(err)

	return a1
}

// Must2 panic if err is not nill，or return 2 result
func Must2[T1, T2 any](a1 T1, a2 T2, err error) (T1, T2) {
	mustCheckError(err)

	return a1, a2
}

// Must3 panic if err is not nill，or return 3 result
func Must3[T1, T2, T3 any](a1 T1, a2 T2, a3 T3, err error) (T1, T2, T3) {
	mustCheckError(err)

	return a1, a2, a3
}

// Must4 panic if err is not nill，or return 4 result
func Must4[T1, T2, T3, T4 any](a1 T1, a2 T2, a3 T3, a4 T4, err error) (T1, T2, T3, T4) {
	mustCheckError(err)

	return a1, a2, a3, a4
}

// Must5 panic if err is not nill，or return 5 result
func Must5[T1, T2, T3, T4, T5 any](a1 T1, a2 T2, a3 T3, a4 T4, a5 T5, err error) (T1, T2, T3, T4, T5) {
	mustCheckError(err)

	return a1, a2, a3, a4, a5
}

func mustCheckError(err error) {
	if err != nil {
		panic(NewError(efrom(err)))
	}
}

var _ = mustCheckArgs

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
