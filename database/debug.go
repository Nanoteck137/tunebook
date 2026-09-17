package database

import (
	"fmt"

	"github.com/nanoteck137/tunebook/tools/pretty"
)

func DebugSQL(query Query) string {
	sql, params, err := query.ToSQL()
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	return pretty.PrettyPrintSQL(sql, params)
}
