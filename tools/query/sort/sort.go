package sort

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nanoteck137/tunebook/tools/query"
)

type Sort struct {
	Orderings []query.Ordering
}

func Parse(input string) (*Sort, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return &Sort{}, nil
	}

	// Check for special modes
	switch strings.ToLower(input) {
	case "random":
		return &Sort{
			Orderings: []query.Ordering{&query.RandomOrdering{}},
		}, nil
	case "score":
		return &Sort{
			Orderings: []query.Ordering{&query.ScoreOrdering{}},
		}, nil
	case "recent":
		return &Sort{
			Orderings: []query.Ordering{
				&query.FieldOrdering{
					Field: &query.Field{Name: "created"},
					Dir:   query.DirDesc,
				},
			},
		}, nil
	}

	// Check for shuffle mode
	if strings.HasPrefix(strings.ToLower(input), "shuffle=") {
		seedStr := strings.TrimPrefix(input, "shuffle=")
		seedStr = strings.TrimPrefix(seedStr, "SHUFFLE=")
		seed, err := strconv.ParseInt(seedStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid shuffle seed: %s", seedStr)
		}
		return &Sort{
			Orderings: []query.Ordering{&query.ShuffleOrdering{Seed: seed}},
		}, nil
	}

	// Parse field orderings
	return parseFieldOrderings(input)
}

func parseFieldOrderings(input string) (*Sort, error) {
	parts := strings.Split(input, ",")
	orderings := make([]query.Ordering, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		ordering, err := parseSingleOrdering(part)
		if err != nil {
			return nil, err
		}
		orderings = append(orderings, ordering)
	}

	if len(orderings) == 0 {
		return nil, fmt.Errorf("no valid orderings found")
	}

	return &Sort{Orderings: orderings}, nil
}

func parseSingleOrdering(part string) (query.Ordering, error) {
	// Extract null ordering if present
	nullOrder := query.NullOrderingDefault
	lowerPart := strings.ToLower(part)
	
	if idx := strings.Index(lowerPart, " nulls first"); idx != -1 {
		nullOrder = query.NullOrderingFirst
		part = strings.TrimSpace(part[:idx])
		lowerPart = strings.ToLower(part)
	} else if idx := strings.Index(lowerPart, " nulls last"); idx != -1 {
		nullOrder = query.NullOrderingLast
		part = strings.TrimSpace(part[:idx])
		lowerPart = strings.ToLower(part)
	}

	// Check for +/- prefix syntax
	if strings.HasPrefix(part, "+") {
		fieldName := strings.TrimSpace(strings.TrimPrefix(part, "+"))
		if fieldName == "" {
			return nil, fmt.Errorf("empty field name after '+'")
		}
		return &query.FieldOrdering{
			Field:     &query.Field{Name: fieldName},
			Dir:       query.DirAsc,
			NullOrder: nullOrder,
		}, nil
	}

	if strings.HasPrefix(part, "-") {
		fieldName := strings.TrimSpace(strings.TrimPrefix(part, "-"))
		if fieldName == "" {
			return nil, fmt.Errorf("empty field name after '-'")
		}
		return &query.FieldOrdering{
			Field:     &query.Field{Name: fieldName},
			Dir:       query.DirDesc,
			NullOrder: nullOrder,
		}, nil
	}

	// Check for "field asc/desc" syntax
	if strings.HasSuffix(lowerPart, " asc") {
		fieldName := strings.TrimSpace(part[:len(part)-4])
		if fieldName == "" {
			return nil, fmt.Errorf("empty field name before 'asc'")
		}
		return &query.FieldOrdering{
			Field:     &query.Field{Name: fieldName},
			Dir:       query.DirAsc,
			NullOrder: nullOrder,
		}, nil
	}

	if strings.HasSuffix(lowerPart, " desc") {
		fieldName := strings.TrimSpace(part[:len(part)-5])
		if fieldName == "" {
			return nil, fmt.Errorf("empty field name before 'desc'")
		}
		return &query.FieldOrdering{
			Field:     &query.Field{Name: fieldName},
			Dir:       query.DirDesc,
			NullOrder: nullOrder,
		}, nil
	}

	// Check if the entire part is just "asc" or "desc" (with or without whitespace)
	trimmed := strings.TrimSpace(part)
	if strings.ToLower(trimmed) == "asc" || strings.ToLower(trimmed) == "desc" {
		return nil, fmt.Errorf("empty field name before '%s'", strings.ToLower(trimmed))
	}

	// Default to ascending if no direction specified
	return &query.FieldOrdering{
		Field:     &query.Field{Name: part},
		Dir:       query.DirAsc,
		NullOrder: nullOrder,
	}, nil
}
