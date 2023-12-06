package do

import (
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "int",
			args: args{
				s: []int{1, 22, 33, 22, 33, 4},
			},
			want: []int{1, 22, 33, 4},
		},
		{
			name: "int-2",
			args: args{
				s: []int{2, 23, 33, 22, 33, 4},
			},
			want: []int{2, 23, 33, 22, 4},
		},
		{
			name: "int-3",
			args: args{
				s: []int{2, 23, 33, 22, 4},
			},
			want: []int{2, 23, 33, 22, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIn(t *testing.T) {
	type args struct {
		s []int
		e int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "not-in",
			args: args{
				s: []int{1, 2, 3},
				e: 0,
			},
			want: false,
		},
		{
			name: "in",
			args: args{
				s: []int{1, 2, 3},
				e: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := In(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}
