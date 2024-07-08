package do

import (
	"fmt"
	"log"
	"runtime/debug"
)

// Go run f in a new goroutine with defer recover
func Go[P any](c C, p P, f func(C, P)) {
	GoR(c, p, func(c C, p P) struct{} {
		f(c, p)
		return struct{}{}
	})
}

// GoR run f in a new goroutine with defer recover, and return a chan of R
func GoR[P, R any](c C, p P, f func(C, P) R) <-chan R {
	ch := make(chan R, 1)

	go func(c C, p P) {
		r := CallInDefRec(c, p, f)

		ch <- r
	}(c, p)

	return ch
}

// CallInDefRec run f with defer recover to catch panic
func CallInDefRec[P, R any](c C, p P, f func(C, P) R) R {
	defer func() {
		if v := recover(); v != nil {
			log.Printf("panic: %v \nstack: %s", v, debug.Stack())
		}
	}()

	return f(c, p)
}

// CallInDefRec2 run f with defer recover to catch panic
func CallInDefRec2[P, R any](c C, p P, f func(C, P) (R, E)) (r R, e E) {
	defer func() {
		if v := recover(); v != nil {
			e = fmt.Errorf("failed: %v", v)

			log.Printf("panic: %v \nstack: %s", v, debug.Stack())
		}
	}()

	return f(c, p)
}
