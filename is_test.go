package do

import "testing"

func TestIsZero(t *testing.T) {
	type args[T comparable] struct {
		v T
	}
	type tcase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []tcase[int]{
		// TODO: Add test cases.
		{
			name: "int is zero",
			args: args[int]{
				v: 0,
			},
			want: true,
		},
		{
			name: "int not zero",
			args: args[int]{
				v: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsZero(tt.args.v); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}

	{
		var a = 1
		tests := []tcase[*int]{
			// TODO: Add test cases.
			{
				name: "pointer of int is zero",
				args: args[*int]{
					v: nil,
				},
				want: true,
			},
			{
				name: "pointer of int not zero",
				args: args[*int]{
					v: &a,
				},
				want: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := IsZero(tt.args.v); got != tt.want {
					t.Errorf("IsZero() = %v, want %v", got, tt.want)
				}
			})
		}
	}

	{
		tests := []tcase[string]{
			// TODO: Add test cases.
			{
				name: "string is zero",
				args: args[string]{
					v: "",
				},
				want: true,
			},
			{
				name: "string not zero",
				args: args[string]{
					v: "a",
				},
				want: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := IsZero(tt.args.v); got != tt.want {
					t.Errorf("IsZero() = %v, want %v", got, tt.want)
				}
			})
		}
	}

	{
		tests := []tcase[float64]{
			// TODO: Add test cases.
			{
				name: "float64 is zero",
				args: args[float64]{
					v: 0.0,
				},
				want: true,
			},
			{
				name: "float64 not zero",
				args: args[float64]{
					v: 1.0,
				},
				want: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := IsZero(tt.args.v); got != tt.want {
					t.Errorf("IsZero() = %v, want %v", got, tt.want)
				}
			})
		}
	}

	{
		tests := []tcase[bool]{
			// TODO: Add test cases.
			{
				name: "bool is zero",
				args: args[bool]{
					v: false,
				},
				want: true,
			},
			{
				name: "bool not zero",
				args: args[bool]{
					v: true,
				},
				want: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := IsZero(tt.args.v); got != tt.want {
					t.Errorf("IsZero() = %v, want %v", got, tt.want)
				}
			})
		}
	}

	// error和[]int都不满足comparable

	{
		type cc struct {
			name string
		}
		var c cc
		tests := []tcase[cc]{
			// TODO: Add test cases.
			{
				name: "struct literal is zero",
				args: args[cc]{
					v: cc{},
				},
				want: true,
			},
			{
				name: "struct is zero",
				args: args[cc]{
					v: c,
				},
				want: true,
			},
			{
				name: "struct not zero",
				args: args[cc]{
					v: cc{"jd"},
				},
				want: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := IsZero(tt.args.v); got != tt.want {
					t.Errorf("IsZero() = %v, want %v", got, tt.want)
				}
			})
		}

		{
			type cc struct {
				name string
			}
			tests := []tcase[*cc]{
				// TODO: Add test cases.
				{
					name: "pointer of struct is zero",
					args: args[*cc]{
						v: nil,
					},
					want: true,
				},
				{
					name: "pointer of struct not zero",
					args: args[*cc]{
						v: &cc{"jd"},
					},
					want: false,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					if got := IsZero(tt.args.v); got != tt.want {
						t.Errorf("IsZero() = %v, want %v", got, tt.want)
					}
				})
			}
		}
	}
}
