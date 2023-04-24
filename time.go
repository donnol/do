package do

import "time"

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
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
}

func ThisMonthFirst() time.Time {
	now := time.Now()
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
}

func ThisYearFirst() time.Time {
	now := time.Now()
	y, _, _ := now.Date()
	return time.Date(y, 1, 1, 0, 0, 0, 0, now.Location())
}
