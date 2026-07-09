package database

import (
	"fmt"

	"github.com/nanoteck137/tunebook/dev"
)

func DebugSQL(query Query) string {
	sql, params, err := query.ToSQL()
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	return dev.PrettyPrintSQL(sql, params)
}
