package do

import "testing"

func TestEqual(t *testing.T) {
	type args struct {
		l int
		r int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equal",
			args: args{
				l: 1,
				r: 1,
			},
			want: true,
		},
		{
			name: "not equal",
			args: args{
				l: 1,
				r: 2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.l, tt.args.r); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
