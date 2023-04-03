package do

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryWithTimes(t *testing.T) {
	ctx := context.Background()
	if err := RetryWithTimes(ctx, 2, func(ctx context.Context) (canRetry bool, err error) {
		return true, errors.New("retry please")
	}); err != nil && err.Error() != "retry please" {
		t.Error(err)
	}
}

func TestRetryWithDeadline(t *testing.T) {
	ctx := context.Background()
	if err := RetryWithDeadline(ctx, time.Now().Add(3*time.Second), func(ctx context.Context) (canRetry bool, err error) {
		return true, errors.New("retry please")
	}); err != nil && err.Error() != "retry please" {
		t.Error(err)
	}
}
