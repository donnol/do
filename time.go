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
