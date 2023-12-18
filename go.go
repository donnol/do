package do

import (
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
		r := callInDefRec(c, p, f)

		ch <- r
	}(c, p)

	return ch
}

// callInDefRec run f with defer recover to catch panic
func callInDefRec[P, R any](c C, p P, f func(C, P) R) R {
	defer func() {
		if v := recover(); v != nil {
			log.Printf("panic stack: %s", debug.Stack())
		}
	}()

	return f(c, p)
}
