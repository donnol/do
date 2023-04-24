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

func TestTodayZero(t *testing.T) {
	today := TodayZero()
	if today.Hour() != 0 {
		t.Errorf("bad day first hour is %v", today.Hour())
	}
	if today.Minute() != 0 {
		t.Errorf("bad day first minute is %v", today.Minute())
	}
	if today.Second() != 0 {
		t.Errorf("bad day first second is %v", today.Second())
	}
	thisMonthFirst := ThisMonthFirst()
	if thisMonthFirst.Day() != 1 {
		t.Errorf("bad month first day is %v", thisMonthFirst.Day())
	}
	if today.Month() != thisMonthFirst.Month() {
		t.Errorf("bad month: %v != %v", today.Month(), thisMonthFirst.Month())
	}
	thisYearFirst := ThisYearFirst()
	if thisYearFirst.Month() != 1 {
		t.Errorf("bad year first month is %v", thisYearFirst.Month())
	}
	if thisYearFirst.Day() != 1 {
		t.Errorf("bad year first day is %v", thisYearFirst.Day())
	}
	if thisYearFirst.Year() != today.Year() {
		t.Errorf("bad year: %v != %v", thisYearFirst.Year(), today.Year())
	}
	if thisYearFirst.Year() != thisMonthFirst.Year() {
		t.Errorf("bad year: %v != %v", thisYearFirst.Year(), thisMonthFirst.Year())
	}
}
