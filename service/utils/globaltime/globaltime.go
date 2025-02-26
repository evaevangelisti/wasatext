package globaltime

import "time"

var FixedTime time.Time

func Now() time.Time {
	if !FixedTime.IsZero() {
		return FixedTime
	}

	return time.Now()
}

func Since(startTime time.Time) time.Duration {
	return Now().Sub(startTime)
}
