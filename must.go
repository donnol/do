package do

import (
	"fmt"
	"log"
)

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

// Log log the err if err is not nill and continue
func Log(err error) {
	logError(err)
}

// Log1 log the err if err is not nill，and continue with 1 result
func Log1[T any](a1 T, err error) T {
	logError(err)

	return a1
}

// Log2 log the err if err is not nill，and continue with 2 result
func Log2[T1, T2 any](a1 T1, a2 T2, err error) (T1, T2) {
	logError(err)

	return a1, a2
}

// Log3 log the err if err is not nill，and continue with 3 result
func Log3[T1, T2, T3 any](a1 T1, a2 T2, a3 T3, err error) (T1, T2, T3) {
	logError(err)

	return a1, a2, a3
}

// Log4 log the err if err is not nill，and continue with 4 result
func Log4[T1, T2, T3, T4 any](a1 T1, a2 T2, a3 T3, a4 T4, err error) (T1, T2, T3, T4) {
	logError(err)

	return a1, a2, a3, a4
}

// Log5 log the err if err is not nill，and continue with 5 result
func Log5[T1, T2, T3, T4, T5 any](a1 T1, a2 T2, a3 T3, a4 T4, a5 T5, err error) (T1, T2, T3, T4, T5) {
	logError(err)

	return a1, a2, a3, a4, a5
}

func logError(err error) {
	if err != nil {
		log.Println(NewError(efrom(err)))
	}
}

// Ignore1 ignore the last variable，and continue with 1 result
func Ignore1[T1, T2 any](a1 T1, i T2) T1 {
	return a1
}

// Ignore2 ignore the last variable，and continue with 2 result
func Ignore2[T1, T2, T3 any](a1 T1, a2 T2, i T3) (T1, T2) {
	return a1, a2
}

// Ignore3 ignore the last variable，and continue with 3 result
func Ignore3[T1, T2, T3, T4 any](a1 T1, a2 T2, a3 T3, i T4) (T1, T2, T3) {
	return a1, a2, a3
}

// Ignore4 ignore the last variable，and continue with 4 result
func Ignore4[T1, T2, T3, T4, T5 any](a1 T1, a2 T2, a3 T3, a4 T4, i T5) (T1, T2, T3, T4) {
	return a1, a2, a3, a4
}

// Ignore5 ignore the last variable，and continue with 5 result
func Ignore5[T1, T2, T3, T4, T5, T6 any](a1 T1, a2 T2, a3 T3, a4 T4, a5 T5, i T6) (T1, T2, T3, T4, T5) {
	return a1, a2, a3, a4, a5
}
