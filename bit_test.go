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
