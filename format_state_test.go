package do

import (
	"testing"
)

func Test_checkPrintf(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name       string
		args       args
		wantOk     bool
		wantArgNum int
	}{
		// TODO: Add test cases.
		// TODO: Add test cases.
		{
			name: "0",
			args: args{
				format: "I am",
			},
			wantOk:     false,
			wantArgNum: 0,
		},
		{
			name: "%%",
			args: args{
				format: "I am %%",
			},
			wantOk:     true,
			wantArgNum: 0,
		},
		{
			name: "1",
			args: args{
				format: "I am %s",
			},
			wantOk:     true,
			wantArgNum: 1,
		},
		{
			name: "*s",
			args: args{
				format: "I am %*s",
			},
			wantOk:     true,
			wantArgNum: 2,
		},
		{
			name: "2",
			args: args{
				format: "I am %s%s",
			},
			wantOk:     true,
			wantArgNum: 2,
		},
		{
			name: "v",
			args: args{
				format: "I am %v",
			},
			wantOk:     true,
			wantArgNum: 1,
		},
		{
			name: "#",
			args: args{
				format: "I am %#v",
			},
			wantOk:     true,
			wantArgNum: 1,
		},
		{
			name: "#2",
			args: args{
				format: "I am %#v %#v",
			},
			wantOk:     true,
			wantArgNum: 2,
		},
		{
			name: "d",
			args: args{
				format: "I am %2d %1d %+d %-d",
			},
			wantOk:     true,
			wantArgNum: 4,
		},
		{
			name: "f",
			args: args{
				format: "I am %2.0f %0.1f %+f %-f %%",
			},
			wantOk:     true,
			wantArgNum: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotArgNum := checkPrintf(tt.args.format)
			if gotOk != tt.wantOk {
				t.Errorf("checkPrintf() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotArgNum != tt.wantArgNum {
				t.Errorf("checkPrintf() gotArgNum = %v, want %v", gotArgNum, tt.wantArgNum)
			}
		})
	}
}
