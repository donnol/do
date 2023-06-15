package do

import (
	"reflect"
	"testing"
)

type M struct {
	name string
}

func TestPtr(t *testing.T) {
	type args[E any] struct {
		e E
		p *E
	}
	type Case[E any] struct {
		name string
		args args[E]
	}
	tests := []Case[int]{
		// TODO: Add test cases.
		{
			name: "int",
			args: func() args[int] {
				v, p := ptrCase()
				arg := args[int]{
					e: v,
					p: p,
				}
				return arg
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPtr := PtrTo(tt.args.e); !reflect.DeepEqual(gotPtr, tt.args.p) {
				t.Errorf("PtrTo() = %v, want %v", gotPtr, tt.args.p)
			}
		})
	}

	{
		tests := []Case[M]{
			// TODO: Add test cases.
			{
				name: "M",
				args: func() args[M] {
					v := M{"go"}
					p := &v
					arg := args[M]{
						e: v,
						p: p,
					}
					return arg
				}(),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotPtr := PtrTo(tt.args.e); !reflect.DeepEqual(gotPtr, tt.args.p) {
					t.Errorf("PtrTo() = %v, want %v", gotPtr, tt.args.p)
				}
			})
		}
	}

	{
		tests := []Case[any]{
			// TODO: Add test cases.
			{
				name: "any",
				args: func() args[any] {
					var v any
					p := &v
					arg := args[any]{
						e: v,
						p: p,
					}
					return arg
				}(),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotPtr := PtrTo(tt.args.e); !reflect.DeepEqual(gotPtr, tt.args.p) {
					t.Errorf("PtrTo() = %v, want %v", gotPtr, tt.args.p)
				}
			})
		}
	}

	{
		tests := []Case[I]{
			// TODO: Add test cases.
			{
				name: "interface",
				args: func() args[I] {
					var v I
					p := &v
					arg := args[I]{
						e: v,
						p: p,
					}
					return arg
				}(),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotPtr := PtrTo(tt.args.e); !reflect.DeepEqual(gotPtr, tt.args.p) {
					t.Errorf("PtrTo() = %v, want %v", gotPtr, tt.args.p)
				}
			})
		}
	}

	{
		p := PtrTo(any(nil))
		if p == nil {
			t.Errorf("bad case, p is nil")
		}
		if p != nil && *p != nil {
			t.Errorf("*p is not nil: %v", *p)
		}
	}
}

func ptrCase() (int, *int) {
	var a int = 1
	return a, &a
}

type I interface {
	String()
}

func TestZero(t *testing.T) {
	if gotE := Zero[int](); !reflect.DeepEqual(gotE, 0) {
		t.Errorf("Zero() = %v, want %v", gotE, 0)
	}

	if gotE := Zero[string](); !reflect.DeepEqual(gotE, "") {
		t.Errorf("Zero() = %v, want %v", gotE, "")
	}

	if gotE := Zero[float64](); !reflect.DeepEqual(gotE, 0.0) {
		t.Errorf("Zero() = %v, want %v", gotE, 0.0)
	}

	if gotE := Zero[M](); !reflect.DeepEqual(gotE, M{}) {
		t.Errorf("Zero() = %v, want %v", gotE, M{})
	}

	if gotE := Zero[[]M](); !reflect.DeepEqual(gotE, []M(nil)) {
		t.Errorf("Zero() = %v, want %v", gotE, []M(nil))
	}

	if gotE := Zero[map[int]M](); !reflect.DeepEqual(gotE, map[int]M(nil)) {
		t.Errorf("Zero() = %v, want %v", gotE, map[int]M(nil))
	}

	if gotE := Zero[chan M](); !reflect.DeepEqual(gotE, chan M(nil)) {
		t.Errorf("Zero() = %v, want %v", gotE, chan M(nil))
	}

	if gotE := Zero[any](); !reflect.DeepEqual(gotE, nil) {
		t.Errorf("Zero() = %v, want %v", gotE, nil)
	}

	if gotE := Zero[I](); !reflect.DeepEqual(gotE, nil) {
		t.Errorf("Zero() = %v, want %v", gotE, nil)
	}

	if gotE := Zero[func()](); !reflect.DeepEqual(gotE, (func())(nil)) {
		t.Errorf("Zero() = %p, want %p", gotE, (func())(nil))
	}
}
