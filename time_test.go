package do

import (
	"testing"
	"time"
)

func TestIsExpired(t *testing.T) {
	type args struct {
		deadline time.Time
		now      time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "zero",
			args: args{
				deadline: time.Time{},
				now:      time.Date(2023, 04, 24, 01, 02, 03, 00, time.Local),
			},
			want: false,
		},
		{
			name: "equal",
			args: args{
				deadline: time.Date(2023, 04, 24, 01, 02, 03, 00, time.Local),
				now:      time.Date(2023, 04, 24, 01, 02, 03, 00, time.Local),
			},
			want: true,
		},
		{
			name: "expired",
			args: args{
				deadline: time.Date(2023, 04, 24, 00, 02, 03, 00, time.Local),
				now:      time.Date(2023, 04, 24, 01, 02, 03, 00, time.Local),
			},
			want: true,
		},
		{
			name: "not expired",
			args: args{
				deadline: time.Date(2023, 04, 24, 02, 02, 03, 00, time.Local),
				now:      time.Date(2023, 04, 24, 01, 02, 03, 00, time.Local),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExpired(tt.args.deadline, tt.args.now); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
