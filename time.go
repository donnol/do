package do

import (
	"time"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
)

var (
	Location = time.FixedZone("CST", 8*3600) // 东八，Asia/Shanghai
)

func Date(year int, month time.Month, day int, loc *time.Location) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func DateLocal(year int, month time.Month, day int) time.Time {
	return Date(year, month, day, time.Local)
}

// IsExpired show if deadline is expired compared to now
// always return false if deadline is zero
func IsExpired(deadline, now time.Time) bool {
	if deadline.IsZero() {
		return false
	}
	return !deadline.After(now)
}

func TodayZero() time.Time {
	now := time.Now()
	return DayZero(now)
}

func ThisMonthFirst() time.Time {
	now := time.Now()
	return MonthFirst(now)
}

func ThisYearFirst() time.Time {
	now := time.Now()
	return YearFirst(now)
}

func DayZero(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func MonthFirst(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

func YearFirst(t time.Time) time.Time {
	y, _, _ := t.Date()
	return time.Date(y, 1, 1, 0, 0, 0, 0, t.Location())
}

// ParseTime parse time string t with layout s one by one; if layouts is empty, it will use "2006-01-02 15:04:05" as default
func ParseTime(t string, layouts ...string) (r time.Time, err error) {
	if len(layouts) == 0 {
		layouts = append(layouts, "2006-01-02 15:04:05")
	}
	for _, layout := range layouts {
		r, err = time.ParseInLocation(layout, t, time.Local)
		if err == nil {
			return
		}
	}
	return
}

// AgeByBirth get age and unit by birthday. It will use time.Now() if not input now.
// If exist year, ignore month and day; If exist month, ignore day.
func AgeByBirth(birthday time.Time, now ...time.Time) (age int, unit string) {
	var nt time.Time
	if len(now) == 0 {
		nt = time.Now()
	} else {
		nt = now[0]
	}

	// 拿到时间间隔
	bd := nt.Sub(birthday)
	if bd < 0 {
		return
	}
	if bd == 0 {
		age = 1
		unit = "天"
	}

	_, bdm, bdd := birthday.Date()
	_, im, id := nt.Date()

	zero := time.Time{}
	ny, nm, nd := zero.Date() // 1 1 1

	y, m, d := zero.Add(bd).Date()
	ry, rm, rd := y-ny, m-nm, d-nd

	if ry > 0 {
		age = ry
		unit = "岁"

		// 如果月份和日子为0时，两个日期的月份和日子中的一个不相等，则说明未足年
		if (rm == 0 && rd == 0) && (im != bdm || id != bdd) {
			age--
		}

	} else if rm > 0 {
		age = int(rm)
		unit = "月"
	} else if rd > 0 {
		age = rd
		unit = "天"
	}

	return
}
