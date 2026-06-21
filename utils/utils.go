package utils

import (
	"database/sql"
	"errors"
	"math"
	"os"
	"strings"
	"time"
)

func TotalPages(perPage, totalItems int) int {
	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}

func SplitTagString(s string) []string {
	tags := []string{}
	if s != "" {
		tags = strings.Split(s, ",")
	}

	return tags
}

func SqlNullToStringPtr(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}

	return nil
}

func StringToPtr(value string) *string {
	if value != "" {
		return &value
	}

	return nil
}

func SqlNullToInt64Ptr(value sql.NullInt64) *int64 {
	if value.Valid {
		return &value.Int64
	}

	return nil
}

func SqlNullToFloat64Ptr(value sql.NullFloat64) *float64 {
	if value.Valid {
		return &value.Float64
	}

	return nil
}

func CreateDirectories(dirs []string) error {
	for _, dir := range dirs {
		err := os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	return nil
}

func Pointer[T any](val T) *T {
	return &val
}

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
