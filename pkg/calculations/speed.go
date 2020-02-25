package calculations

import "time"

// Speed in miles per hours based on timestamp deltas and given distance in miles.
func SpeedMilesPerHour(timestamp1, timestamp2 int64, distanceMiles float64) int {
	var d time.Duration
	t1 := time.Unix(timestamp1, 0)
	t2 := time.Unix(timestamp2, 0)

	if t1.Before(t2) {
		d = t2.Sub(t1)
	} else {
		d = t1.Sub(t2)
	}

	if d.Hours() > 0 {
		speed := distanceMiles / d.Hours()
		return int(speed)
	}

	return 0
}
