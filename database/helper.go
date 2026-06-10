package database

import (
	"math"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nrednav/cuid2"
)

func addToRecord[T any](record goqu.Record, name string, change types.Change[T]) {
	if change.Changed {
		record[name] = change.Value
	}
}

func createIdGenerator(length int) func() string {
	res, err := cuid2.Init(cuid2.WithLength(length))
	if err != nil {
		panic(err)
	}

	return res
}

func totalPages(perPage, totalItems int) int {
	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}
