package do

import (
	"reflect"
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

func TestParseTime(t *testing.T) {
	type args struct {
		t       string
		layouts []string
	}
	tests := []struct {
		name    string
		args    args
		wantR   time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "empty layout",
			args: args{
				t:       "2011-02-04 23:33:10",
				layouts: []string{},
			},
			wantR:   time.Date(2011, 02, 04, 23, 33, 10, 0, time.Local),
			wantErr: false,
		},
		{
			name: "one layout",
			args: args{
				t: "2011-02-04T23:33:10Z",
				layouts: []string{
					"2006-01-02T15:04:05Z",
				},
			},
			wantR:   time.Date(2011, 02, 04, 23, 33, 10, 0, time.Local),
			wantErr: false,
		},
		{
			name: "many layout",
			args: args{
				t: "2011-02-04T23:33:10Z",
				layouts: []string{
					"2006-01-02 15:04:05 ",
					"2006-01-02T15:04:05Z",
				},
			},
			wantR:   time.Date(2011, 02, 04, 23, 33, 10, 0, time.Local),
			wantErr: false,
		},
		{
			name: "many layout but failed",
			args: args{
				t: "2011-02-04T23:33:10Z",
				layouts: []string{
					"2006-01-02 15:04:05 ",
					"2006-01-02",
				},
			},
			wantR:   time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := ParseTime(tt.args.t, tt.args.layouts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("ParseTime() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestNowTwice(t *testing.T) {
	n1 := time.Now()
	n2 := time.Now()
	_, _ = n1, n2

	// 大概率相等，也会有不相等的出现
	// Assert(t, !n2.Before(n1), true)
	// Assert(t, n1.UnixMilli(), n2.UnixMilli())
	// Assert(t, n1.UnixMicro(), n2.UnixMicro())
	// Assert(t, n1.String(), n2.String())
}

func TestAgeByBirth(t *testing.T) {
	type args struct {
		birthday time.Time
		now      []time.Time
	}
	tests := []struct {
		name     string
		args     args
		wantAge  int
		wantUnit string
	}{
		// TODO: Add test cases.
		{
			name: "2y",
			args: args{
				birthday: time.Date(2022, 01, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 01, 22, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  2,
			wantUnit: "岁",
		},
		{
			name: "2y",
			args: args{
				birthday: time.Date(2022, 01, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 01, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  1,
			wantUnit: "岁",
		},
		{
			name: "20y-01",
			args: args{
				birthday: time.Date(2002, 01, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 01, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-02",
			args: args{
				birthday: time.Date(2002, 02, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 02, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-03",
			args: args{
				birthday: time.Date(2002, 03, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 03, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-04",
			args: args{
				birthday: time.Date(2002, 04, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 04, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-05",
			args: args{
				birthday: time.Date(2002, 05, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 05, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-06",
			args: args{
				birthday: time.Date(2002, 06, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 06, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-07",
			args: args{
				birthday: time.Date(2002, 07, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 07, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-08",
			args: args{
				birthday: time.Date(2002, 8, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 8, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-09",
			args: args{
				birthday: time.Date(2002, 9, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 9, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-10",
			args: args{
				birthday: time.Date(2002, 10, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 10, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-11",
			args: args{
				birthday: time.Date(2002, 11, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 11, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "20y-12",
			args: args{
				birthday: time.Date(2002, 12, 22, 0, 0, 0, 0, time.Local),
				now: []time.Time{
					time.Date(2024, 12, 21, 0, 0, 0, 0, time.Local),
				},
			},
			wantAge:  21,
			wantUnit: "岁",
		},
		{
			name: "1y",
			args: args{
				birthday: time.Date(2023, 01, 22, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 1, 22, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  1,
			wantUnit: "岁",
		},
		{
			name: "未足年",
			args: args{
				birthday: time.Date(2023, 01, 22, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 1, 21, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  11,
			wantUnit: "月",
		},
		{
			name: "未足年",
			args: args{
				birthday: time.Date(2023, 02, 22, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 1, 22, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  11,
			wantUnit: "月",
		},
		{
			name: "1m",
			args: args{
				birthday: time.Date(2023, 12, 22, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 1, 22, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  1,
			wantUnit: "月",
		},
		{
			name: "1m",
			args: args{
				birthday: time.Date(2023, 11, 22, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2023, 12, 22, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  30,
			wantUnit: "天",
		},
		{
			name: "1d",
			args: args{
				birthday: time.Date(2024, 01, 21, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 1, 22, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  1,
			wantUnit: "天",
		},
		{
			name: "1d",
			args: args{
				birthday: time.Date(2024, 01, 22, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 1, 22, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  1,
			wantUnit: "天",
		},
		{
			name: "rev",
			args: args{
				birthday: time.Date(2024, 01, 22, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 1, 21, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  0,
			wantUnit: "",
		},
		{
			name: "before-dates",
			args: args{
				// 2006年8月4日
				birthday: time.Date(2006, 8, 4, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 04, 30, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  17,
			wantUnit: "岁",
		},
		{
			name: "before-2d",
			args: args{
				// 2006年8月4日
				birthday: time.Date(2006, 8, 4, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 8, 2, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  17,
			wantUnit: "岁",
		},
		{
			name: "before-1d",
			args: args{
				// 2006年8月4日
				birthday: time.Date(2006, 8, 4, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 8, 3, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  17,
			wantUnit: "岁",
		},
		{
			name: "equal-date",
			args: args{
				// 2006年8月4日
				birthday: time.Date(2006, 8, 4, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 8, 4, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  18,
			wantUnit: "岁",
		},
		{
			name: "after-1d",
			args: args{
				// 2006年8月4日
				birthday: time.Date(2006, 8, 4, 0, 0, 0, 0, time.Local),
				now:      []time.Time{time.Date(2024, 8, 5, 0, 0, 0, 0, time.Local)},
			},
			wantAge:  18,
			wantUnit: "岁",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAge, gotUnit := AgeByBirth(tt.args.birthday, tt.args.now...)
			if gotAge != tt.wantAge {
				t.Errorf("AgeByBirth() gotAge = %v, want %v", gotAge, tt.wantAge)
			}
			if gotUnit != tt.wantUnit {
				t.Errorf("AgeByBirth() gotUnit = %v, want %v", gotUnit, tt.wantUnit)
			}
		})
	}
}
