package do

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNilDoer = errors.New("f is nil")
)

type Doer = func(context.Context) (canRetry bool, err error)

// RetryWithTimes retry f tryTimes times if f failed
func RetryWithTimes(ctx context.Context, tryTimes int, f Doer) error {
	if f == nil {
		return ErrNilDoer
	}

	var err error
	for i := 1; i <= tryTimes; i++ {
		if canRetry, err1 := f(ctx); err1 != nil {
			err = err1

			if !canRetry {
				return err
			}

			time.Sleep(time.Second * (1<<i - 1))
			continue
		}

		err = nil
		break
	}

	return err
}

// RetryWithDeadline retry f before d exceeds if f failed
func RetryWithDeadline(ctx context.Context, d time.Time, f Doer) error {
	if f == nil {
		return ErrNilDoer
	}

	var err error
	var i = 1
	for {
		if canRetry, err1 := f(ctx); err1 != nil {
			err = err1

			if !canRetry {
				return err
			}

			now := time.Now()
			if now.After(d) {
				break
			}

			time.Sleep(time.Second * (1<<i - 1))
			i++
			continue
		}

		err = nil
		break
	}

	return err
}
