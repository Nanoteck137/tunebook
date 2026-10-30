package utils

import "database/sql"

func SqlNullToStringPtr(value sql.NullString) *string {
	if value.Valid {
		return &value.String
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
