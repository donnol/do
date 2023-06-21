package do

import "testing"

func TestIsValidIP(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "ok",
			args: args{
				ip: "127.0.0.1",
			},
			want: true,
		},
		{
			name: "false",
			args: args{
				ip: "127.0.1",
			},
			want: false,
		},
		{
			name: "ipv6",
			args: args{
				ip: "2001:db8::68",
			},
			want: true,
		},
		{
			name: "ipv6 with zone",
			args: args{
				ip: "fe80::1cc0:3e8c:119f:c2e1%ens18",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidIP(tt.args.ip); got != tt.want {
				t.Errorf("IsValidIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
