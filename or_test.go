package do

import (
	"testing"
)

func TestOr(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{nil, 0},
		{[]int{0}, 0},
		{[]int{1}, 1},
		{[]int{0, 2}, 2},
		{[]int{3, 0}, 3},
		{[]int{4, 5}, 4},
		{[]int{0, 6, 7}, 6},
	}
	for _, tc := range cases {
		if got := Or(tc.in...); got != tc.want {
			t.Errorf("cmp.Or(%v) = %v; want %v", tc.in, got, tc.want)
		}
	}
}
