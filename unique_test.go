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

func TestFirst(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name   string
		args   args
		wantT  int
		wantOk bool
	}{
		// TODO: Add test cases.
		{
			name: "0",
			args: args{
				s: []int{},
			},
			wantT:  0,
			wantOk: false,
		},
		{
			name: "1",
			args: args{
				s: []int{1},
			},
			wantT:  1,
			wantOk: true,
		},
		{
			name: "2",
			args: args{
				s: []int{1, 2},
			},
			wantT:  1,
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotOk := First(tt.args.s)
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("First() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotOk != tt.wantOk {
				t.Errorf("First() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestLast(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name   string
		args   args
		wantT  int
		wantOk bool
	}{
		// TODO: Add test cases.
		{
			name: "0",
			args: args{
				s: []int{},
			},
			wantT:  0,
			wantOk: false,
		},
		{
			name: "1",
			args: args{
				s: []int{1},
			},
			wantT:  1,
			wantOk: true,
		},
		{
			name: "2",
			args: args{
				s: []int{1, 2},
			},
			wantT:  2,
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotOk := Last(tt.args.s)
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("First() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotOk != tt.wantOk {
				t.Errorf("First() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	type args struct {
		s []int
		i int
	}
	tests := []struct {
		name   string
		args   args
		wantT  int
		wantOk bool
	}{
		// TODO: Add test cases.
		{
			name: "0",
			args: args{
				s: []int{},
				i: 0,
			},
			wantT:  0,
			wantOk: false,
		},
		{
			name: "1-0",
			args: args{
				s: []int{1},
				i: 0,
			},
			wantT:  1,
			wantOk: true,
		},
		{
			name: "1-1",
			args: args{
				s: []int{1},
				i: 1,
			},
			wantT:  0,
			wantOk: false,
		},
		{
			name: "2-0",
			args: args{
				s: []int{1, 2},
				i: 0,
			},
			wantT:  1,
			wantOk: true,
		},
		{
			name: "2-1",
			args: args{
				s: []int{1, 2},
				i: 1,
			},
			wantT:  2,
			wantOk: true,
		},
		{
			name: "2-3",
			args: args{
				s: []int{1, 2},
				i: 2,
			},
			wantT:  0,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotOk := Index(tt.args.s, tt.args.i)
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("First() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotOk != tt.wantOk {
				t.Errorf("First() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
