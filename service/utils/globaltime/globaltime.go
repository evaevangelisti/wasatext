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

func Parse(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, t)
}

func Format(t time.Time) (string, error) {
	if t.IsZero() {
		return "", nil
	}

	return t.Format(time.RFC3339), nil
}
