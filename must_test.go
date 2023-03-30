package do_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/donnol/do"
)

type CustomError struct {
	code ErrorCode
	msg  string
}

func NewError(
	code ErrorCode,
	msg string,
) error {
	return &CustomError{
		code: code,
		msg:  msg,
	}
}

func (e CustomError) Error() string {
	return fmt.Sprintf("[ERR] code: %d, msg: %s", e.code, e.msg)
}

type ErrorCode int64

const (
	ErrorCodeMust ErrorCode = 1
)

func TestMust(t *testing.T) {
	t.Run("not error", func(t *testing.T) {
		do.Must("not error")
	})
	t.Run("nil", func(t *testing.T) {
		do.Must(nil)
	})
	t.Run("nil error", func(t *testing.T) {
		e := error(nil)
		do.Must(e)
	})
	t.Run("std error come", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if fmt.Sprintf("%v", r) != "raw error: std error come" {
					t.Errorf("recover failed: %v", r)
				}
			}
		}()

		do.Must(errors.New("std error come"))
	})
	t.Run("std error come but empty", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if fmt.Sprintf("%v", r) != "raw error: " {
					t.Errorf("recover failed: %v", r)
				}
			}
		}()

		do.Must(errors.New(""))
	})
	t.Run("custom error come", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if fmt.Sprintf("%v", r) != "raw error: [ERR] code: 1, msg: custom error come" {
					t.Errorf("recover failed: %v", r)
				}
			}
		}()

		do.Must(NewError(1, "custom error come"))
	})
	t.Run("match error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if fmt.Sprintf("%v", r) != "raw error: [ERR] code: 1, msg: custom error come" {
					t.Errorf("recover failed: %v", r)
				}

				match := do.MatchError(r)
				if !match {
					t.Errorf("r is not match with do.Error")
				}
			}
		}()

		do.Must(NewError(1, "custom error come"))
	})
	t.Run("convert error", func(t *testing.T) {
		err := NewError(1, "custom error come")

		defer func() {
			if r := recover(); r != nil {
				if fmt.Sprintf("%v", r) != "raw error: [ERR] code: 1, msg: custom error come" {
					t.Errorf("recover failed: %v", r)
				}

				e, match := do.ConvertError(r)
				if !match {
					t.Errorf("r is not match with do.Error")
				} else {
					ce := e.Inner()
					if ce.Raw() != err {
						t.Errorf("raw error is not equal %v", err)
					}
				}
			}
		}()

		do.Must(err)
	})
}
