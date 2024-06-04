package do

import (
	"reflect"
	"testing"
)

func TestRectangle(t *testing.T) {
	type args struct {
		m       int
		n       int
		initial []int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "0-0",
			args: args{
				m: 0,
				n: 0,
			},
			want: [][]int{},
		},
		{
			name: "4-3",
			args: args{
				m: 4,
				n: 3,
			},
			want: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
		},
		{
			name: "4-3",
			args: args{
				m:       4,
				n:       3,
				initial: []int{-1},
			},
			want: [][]int{
				{-1, -1, -1},
				{-1, -1, -1},
				{-1, -1, -1},
				{-1, -1, -1},
			},
		},
		{
			name: "4-3",
			args: args{
				m:       4,
				n:       3,
				initial: []int{-1, -2},
			},
			want: [][]int{
				{-1, -1, -1},
				{-1, -1, -1},
				{-1, -1, -1},
				{-1, -1, -1},
			},
		},
		{
			name: "4-3",
			args: args{
				m:       4,
				n:       3,
				initial: []int{-1, -2, -3},
			},
			want: [][]int{
				{-1, -2, -3},
				{-1, -2, -3},
				{-1, -2, -3},
				{-1, -2, -3},
			},
		},
		{
			name: "8-6",
			args: args{
				m: 8,
				n: 6,
			},
			want: [][]int{
				{0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0},
			},
		},
		{
			name: "8-6",
			args: args{
				m:       8,
				n:       6,
				initial: []int{1},
			},
			want: [][]int{
				{1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1},
			},
		},
		{
			name: "8-6",
			args: args{
				m:       8,
				n:       6,
				initial: []int{1, 2, 3, 4, 5, 6},
			},
			want: [][]int{
				{1, 2, 3, 4, 5, 6},
				{1, 2, 3, 4, 5, 6},
				{1, 2, 3, 4, 5, 6},
				{1, 2, 3, 4, 5, 6},
				{1, 2, 3, 4, 5, 6},
				{1, 2, 3, 4, 5, 6},
				{1, 2, 3, 4, 5, 6},
				{1, 2, 3, 4, 5, 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rectangle[int](tt.args.m, tt.args.n, tt.args.initial...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rectangle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSquare(t *testing.T) {
	type args struct {
		n       int
		initial []int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "0",
			args: args{
				n: 0,
			},
			want: [][]int{},
		},
		{
			name: "3",
			args: args{
				n: 3,
			},
			want: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
		},
		{
			name: "3",
			args: args{
				n:       3,
				initial: []int{-1},
			},
			want: [][]int{
				{-1, -1, -1},
				{-1, -1, -1},
				{-1, -1, -1},
			},
		},
		{
			name: "3",
			args: args{
				n:       3,
				initial: []int{-1, -2},
			},
			want: [][]int{
				{-1, -1, -1},
				{-1, -1, -1},
				{-1, -1, -1},
			},
		},
		{
			name: "3",
			args: args{
				n:       3,
				initial: []int{-1, -2, -3},
			},
			want: [][]int{
				{-1, -2, -3},
				{-1, -2, -3},
				{-1, -2, -3},
			},
		},
		{
			name: "10",
			args: args{
				n: 10,
			},
			want: [][]int{
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			name: "10",
			args: args{
				n:       10,
				initial: []int{1},
			},
			want: [][]int{
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			},
		},
		{
			name: "10",
			args: args{
				n:       10,
				initial: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: [][]int{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Square[int](tt.args.n, tt.args.initial...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Square() = %v, want %v", got, tt.want)
			}
		})
	}
}
