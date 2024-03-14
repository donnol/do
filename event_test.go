package do

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
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

func TestPipes(t *testing.T) {
	p := "begin"
	before := PipeFrom(func(ctx C, p string) (string, error) {
		return p + " before", nil
	})
	do := PipeFrom(func(ctx C, p string) (string, error) {
		return p + " do", nil
	})
	after := PipeFrom(func(ctx C, p string) (string, error) {
		return p + " after", nil
	})
	// 好笨啊，1.20这都推不出类型；需要用1.22才行：https://go.dev/play/p/C7tpqRGVXDJ
	r := Must1(Pipes[string, string, string, string](context.Background(), p, before, do, after))
	Assert(t, r, "begin before do after")
}

func TestRunEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		buf := new(bytes.Buffer)

		Event(context.Background(), 1, func(ctx C, i int) (string, E) {
			return "h", nil
		}, func(ctx C, s string) (string, error) {
			fmt.Fprintf(buf, "success: %v\n", s)
			return s, nil
		}, func(ctx C, err error) (string, error) {
			fmt.Fprintf(buf, "failed: %v\n", err)
			return "", err
		})

		if !strings.Contains(buf.String(), "success") {
			t.Error("bad case")
		}
	})

	t.Run("failed", func(t *testing.T) {
		buf := new(bytes.Buffer)

		Event(context.Background(), 1, func(ctx C, i int) (string, E) {
			return "", fmt.Errorf("bad case")
		}, func(ctx C, s string) (string, error) {
			fmt.Fprintf(buf, "success: %v\n", s)
			return s, nil
		}, func(ctx C, err error) (string, error) {
			fmt.Fprintf(buf, "failed: %v\n", err)
			return err.Error(), nil
		})

		if !strings.Contains(buf.String(), "failed") {
			t.Error("bad case")
		}
	})
}

func TestEventLoop(t *testing.T) {
	type args struct {
		ctx C
		n   int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				ctx: context.TODO(),
				n:   10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := EventLoop[int, string, string](tt.args.ctx, tt.args.n)

			for _, e := range []struct {
				id int
				s  string
			}{
				{
					id: 1,
					s:  "1",
				},
				{
					id: 2,
					s:  "2",
				},
				{
					id: 3,
					s:  "3",
				},
			} {
				e := e // avoid shadow

				got <- EventEntity[int, string, string]{
					Param: e.id,
					Do: func(ctx context.Context, id int) (string, error) {
						return strconv.Itoa(id), nil
					},
					Success: func(ctx C, id string) (string, error) {
						return id, nil
					},
					Failed: func(ctx C, err error) (string, error) {
						return err.Error(), nil
					},
					Handler: func(id string, err error) {
						if id != e.s {
							t.Errorf("bad case: %s != %s", id, e.s)
						}
					},
				}
			}

			got1 <- struct{}{}
			time.Sleep(time.Millisecond * 200)

			// 停止后继续发送: panic: send on closed channel
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("need panic")
					} else {
						if !strings.Contains(fmt.Sprintf("%v", r), "send on closed channel") {
							t.Errorf("bad case, don't got 'send on closed channel' error: %v", r)
						}
					}
				}()

				for _, e := range []struct {
					id int
					s  string
				}{
					{
						id: 4,
						s:  "4",
					},
				} {
					got <- EventEntity[int, string, string]{
						Param: e.id,
						Do: func(ctx context.Context, id int) (string, error) {
							return strconv.Itoa(id), nil
						},
						Success: func(ctx C, id string) (string, error) {
							return id, nil
						},
						Failed: func(ctx C, err error) (string, error) {
							return err.Error(), nil
						},
						Handler: func(id string, err error) {
							if id != e.s {
								t.Errorf("bad case: %s != %s", id, e.s)
							}
						},
					}
				}
			}()

		})
	}
}
