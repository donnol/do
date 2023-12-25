package do

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestRunWithTimeout(t *testing.T) {
	ctx := context.Background()
	d := time.Second * 2

	// normal
	r, ok := RunWithTimeout(ctx, d, 1, func(ctx context.Context, p int) string {
		return strconv.Itoa(p)
	})
	Assert(t, ok, true)
	Assert(t, r, "1")

	// less than timeout duration
	{
		begin := time.Now()
		r, ok := RunWithTimeout(ctx, d, 1, func(ctx context.Context, p int) string {
			time.Sleep(time.Second * 1)
			return strconv.Itoa(p)
		})
		Assert(t, ok, true)
		Assert(t, r, "1")
		used := time.Since(begin).Seconds()
		Assert(t, used >= 1.0 && used < 1.1, true)
	}

	// timeout
	{
		begin := time.Now()
		r, ok := RunWithTimeout(ctx, d, 1, func(ctx context.Context, p int) string {
			time.Sleep(time.Second * 3)
			return strconv.Itoa(p)
		})
		Assert(t, ok, false)
		Assert(t, r, "")
		used := time.Since(begin).Seconds()
		Assert(t, used >= 2.0 && used < 2.1, true)
	}
}
