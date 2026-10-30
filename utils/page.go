package utils

import "math"

func TotalPages(perPage, totalItems int) int {
	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}
