package do

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type Before struct {
	Name string
}

type Do struct {
	Name string
}

type After struct {
	Name string
}

type Result struct {
	Name string
}

func TestPipe(t *testing.T) {
	type args[B, D, A, R any] struct {
		ctx    context.Context
		b      B
		before PipeFunc[B, D]
		do     PipeFunc[D, A]
		after  PipeFunc[A, R]
	}
	type testCase[B, D, A, R any] struct {
		name    string
		args    args[B, D, A, R]
		wantR   R
		wantErr bool
	}
	tests := []testCase[Before, Do, After, Result]{
		// TODO: Add test cases.
		{
			name: "1",
			args: args[Before, Do, After, Result]{
				ctx: context.Background(),
				b: Before{
					Name: "jd",
				},
				before: func(ctx context.Context, b Before) (r Do, err error) {
					r.Name = b.Name + " is"
					return
				},
				do: func(ctx context.Context, d Do) (r After, err error) {
					r.Name = d.Name + " a"
					return
				},
				after: func(ctx context.Context, a After) (r Result, err error) {
					r.Name = a.Name + " father"
					return
				},
			},
			wantR: Result{Name: "jd is a father"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := Pipe(tt.args.ctx, tt.args.b, tt.args.before, tt.args.do, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Wrap() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestPipeNested(t *testing.T) {
	r, err := Pipe(
		context.Background(),
		1,
		func(ctx C, id int) (string, E) {
			return "1", nil
		},
		func(ctx C, id string) (float64, E) {
			return Pipe(
				ctx,
				id,
				func(ctx C, id string) (float64, E) {
					return 0.2, nil
				},
				func(ctx C, id float64) (string, E) {
					return "h", nil
				},
				func(ctx C, id string) (float64, E) {
					return 2, nil
				},
			)
		},
		func(ctx C, f float64) (uint64, E) {
			return 1, nil
		},
	)
	if err != nil {
		t.Error(err)
	}
	if r != 1 {
		t.Errorf("bad case")
	}
}

func TestRunEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		buf := new(bytes.Buffer)

		Event(context.Background(), 1, func(ctx C, i int) (string, E) {
			return "h", nil
		}, func(s string) string {
			fmt.Fprintf(buf, "success: %v\n", s)
			return s
		}, func(err error) string {
			fmt.Fprintf(buf, "failed: %v\n", err)
			return err.Error()
		})

		if !strings.Contains(buf.String(), "success") {
			t.Error("bad case")
		}
	})

	t.Run("failed", func(t *testing.T) {
		buf := new(bytes.Buffer)

		Event(context.Background(), 1, func(ctx C, i int) (string, E) {
			return "", fmt.Errorf("bad case")
		}, func(s string) string {
			fmt.Fprintf(buf, "success: %v\n", s)
			return s
		}, func(err error) string {
			fmt.Fprintf(buf, "failed: %v\n", err)
			return err.Error()
		})

		if !strings.Contains(buf.String(), "failed") {
			t.Error("bad case")
		}
	})
}
