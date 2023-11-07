package do

import "time"

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
)

var (
	Location = time.FixedZone("CST", 8*3600) // 东八，Asia/Shanghai
)

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
