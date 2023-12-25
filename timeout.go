package do

import (
	"context"
	"time"
)

// RunWithTimeout run func with timeout, return even if the func haven't completed. `ok` will report if the func has completed
func RunWithTimeout[P, R any](ctx C, d time.Duration, p P, f func(ctx C, p P) R) (r R, ok bool) {
	timer := time.NewTimer(d)
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	ch := GoR(ctx, p, f)

	select {
	case <-timer.C:
		return
	case r = <-ch:
		ok = true
		return
	}
}
