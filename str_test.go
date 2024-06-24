package do

import "testing"

func TestFuzzWrap(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				v: "hello",
			},
			want: "%hello%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FuzzWrap(tt.args.v); got != tt.want {
				t.Errorf("FuzzWrap() = %v, want %v", got, tt.want)
			}
		})
	}
}
