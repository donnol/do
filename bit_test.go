package do

import "testing"

func TestIsPowerOf2(t *testing.T) {
	type args struct {
		n uint64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "0",
			args: args{
				n: 0,
			},
			want: false,
		},
		{
			name: "1",
			args: args{
				n: 1,
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				n: 2,
			},
			want: true,
		},
		{
			name: "3",
			args: args{
				n: 3,
			},
			want: false,
		},
		{
			name: "4",
			args: args{
				n: 4,
			},
			want: true,
		},
		{
			name: "5",
			args: args{
				n: 5,
			},
			want: false,
		},
		{
			name: "128",
			args: args{
				n: 128,
			},
			want: true,
		},
		{
			name: "129",
			args: args{
				n: 129,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPowerOf2(tt.args.n); got != tt.want {
				t.Errorf("IsPowerOf2() = %v, want %v", got, tt.want)
			}
		})
	}
}

// From: https://github.com/labuladong/fucking-algorithm/blob/master/%E7%AE%97%E6%B3%95%E6%80%9D%E7%BB%B4%E7%B3%BB%E5%88%97/%E5%B8%B8%E7%94%A8%E7%9A%84%E4%BD%8D%E6%93%8D%E4%BD%9C.md

func TestToLowerByBit(t *testing.T) {
	Assert(t, ('a' | ' '), 'a')
	Assert(t, ('A' | ' '), 'a')
	Assert(t, ('b' & '_'), 'B')
	Assert(t, ('B' & '_'), 'B')
	Assert(t, ('d' ^ ' '), 'D')
	Assert(t, ('D' ^ ' '), 'd')
}

func TestSwapByBit(t *testing.T) {
	var a, b = 1, 2
	a, b = SwapByBit(a, b)

	// 现在 a = 2, b = 1
	Assert(t, a, 2)
	Assert(t, b, 1)
}

func TestAddOneByBit(t *testing.T) {
	var n = 1
	n = -^n

	// 现在 n = 2
	Assert(t, n, 2)
}

func TestSubOneByBit(t *testing.T) {
	var n = 2
	n = ^-n

	// 现在 n = 1
	Assert(t, n, 1)
}

func TestIsSignDiff(t *testing.T) {
	{
		var x, y = -1, 2
		f := IsSignDiff(x, y)
		Assert(t, f, true)
	}

	{
		var x, y = 3, 2
		f := IsSignDiff(x, y)
		Assert(t, f, false)
	}
}

func TestSelf(t *testing.T) {
	// 任何数，异或自身必得0
	for i := -1; i < 10000; i++ {
		//lint:ignore SA4000 this is ok
		r := i ^ i
		Assert(t, r, 0)
	}
}
