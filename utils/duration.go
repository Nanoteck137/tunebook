package utils

import "time"

func PrettyDuration(d time.Duration) string {
	switch {
	case d < time.Millisecond:
		return d.String() // show microseconds as-is
	case d < time.Second:
		return d.Truncate(time.Millisecond).String() // "123ms"
	case d < time.Minute:
		return d.Truncate(time.Second).String() // "42s"
	default:
		return d.Truncate(time.Second).String() // "2h35m42s"
	}
}
