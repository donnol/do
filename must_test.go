package do_test

import (
	"errors"
	"fmt"
	"log"
	"reflect"
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

func fm1() (int, error) {
	return 1, nil
}

func fm2() (int, string, error) {
	return 1, "", nil
}

func fm3() (int, string, string, error) {
	return 1, "", "", nil
}

func fm4() (int, string, string, float64, error) {
	return 1, "", "", 0, nil
}

func fm5() (int, string, string, float64, float64, error) {
	return 1, "", "", 0, 1, nil
}

func TestMustN(t *testing.T) {
	r := do.Must1(fm1())
	if r != 1 {
		t.Errorf("bad case: %v != %v", r, 1)
	}
	{
		r1, r2 := do.Must2(fm2())
		if r1 != 1 && r2 != "" {
			t.Errorf("bad case: %v != %v", r, 1)
		}
	}
	{
		r1, r2, r3 := do.Must3(fm3())
		if r1 != 1 && r2 != "" && r3 != "" {
			t.Errorf("bad case: %v != %v", r, 1)
		}
	}
	{
		r1, r2, r3, r4 := do.Must4(fm4())
		if r1 != 1 && r2 != "" && r3 != "" && r4 != 0 {
			t.Errorf("bad case: %v != %v", r, 1)
		}
	}
	{
		r1, r2, r3, r4, r5 := do.Must5(fm5())
		if r1 != 1 && r2 != "" && r3 != "" && r4 != 0 && r5 != 0 {
			t.Errorf("bad case: %v != %v", r, 1)
		}
	}
}

func TestLog1(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	type args struct {
		a1  int
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "0",
			args: args{
				a1:  0,
				err: errors.New("err show"),
			},
			want: 0,
		},
		{
			name: "1",
			args: args{
				a1:  1,
				err: errors.New("err show"),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := do.Log1(tt.args.a1, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Log1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIgnore(t *testing.T) {
	r := do.Ignore1(1, "2")
	do.Assert(t, r, 1)

	r1, r2 := do.Ignore2(1, "2", 4.0)
	do.Assert(t, r1, 1)
	do.Assert(t, r2, "2")
}

func TestReturn(t *testing.T) {
	sum := func(a, b int) (int, error) {
		return a + b, nil
	}

	// :=
	// 存在新变量，展开后继续是`:=`
	{
		func() error {
			r := do.Return1(sum(1, 2))

			do.Assert(t, r, 3)
			return nil
		}()

		// run `letgo expand` to expand the `do.Return1` into:
		func() error {
			r, err := sum(1, 2)
			if err != nil {
				return err
			}

			do.Assert(t, r, 3)
			return nil
		}()
	}

	// =
	// 不存在新变量，展开后需要添加err变量声明
	{
		func() error {
			var r int
			_ = r

			r = do.Return1(sum(1, 2))

			do.Assert(t, r, 3)
			return nil
		}()

		// run `letgo expand` to expand the `do.Return1` into:
		func() error {
			var r int

			var err error
			r, err = sum(1, 2)
			if err != nil {
				return err
			}

			do.Assert(t, r, 3)
			return nil
		}()
	}
}
