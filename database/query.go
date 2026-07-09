package database

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/parser"
	"github.com/nanoteck137/tunebook/tools/query/planner"
	"github.com/nanoteck137/tunebook/tools/query/schema"
	"github.com/nanoteck137/tunebook/tools/query/sort"
	querysql "github.com/nanoteck137/tunebook/tools/query/sql"
)

// QueryParams contains filter and sort parameters
type QueryParams struct {
	Filter string
	Sort   string
}

// ApplyQuery applies filter and sort to a goqu query using the provided schema
func ApplyQuery(
	q *goqu.SelectDataset,
	s *schema.Schema,
	params QueryParams,
) (*goqu.SelectDataset, error) {
	pl := planner.New(s)
	compiler := querysql.NewCompiler()

	// Parse and plan filter if provided
	if params.Filter != "" {
		p := parser.New(params.Filter)
		expr, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("filter parse error: %w", err)
		}

		plan, err := pl.Plan(expr)
		if err != nil {
			return nil, fmt.Errorf("filter plan error: %w", err)
		}

		result, err := compiler.Compile(plan)
		if err != nil {
			return nil, fmt.Errorf("filter compile error: %w", err)
		}

		if result.Where != nil {
			q = q.Where(result.Where)
		}
	}

	// Parse and apply sort
	var sortOrderings []query.Ordering
	if params.Sort != "" {
		sortObj, err := sort.Parse(params.Sort)
		if err != nil {
			return nil, fmt.Errorf("sort parse error: %w", err)
		}
		sortOrderings = sortObj.Orderings
	}

	// Resolve sort (applies default if empty)
	resolvedOrderings, err := pl.ResolveSort(sortOrderings)
	if err != nil {
		return nil, fmt.Errorf("sort resolution error: %w", err)
	}

	if len(resolvedOrderings) > 0 {
		// Create a temporary plan to compile the orderings
		plan := &query.Plan{OrderBy: resolvedOrderings}
		result, err := compiler.Compile(plan)
		if err != nil {
			return nil, fmt.Errorf("sort compile error: %w", err)
		}

		if len(result.Order) > 0 {
			q = q.Order(result.Order...)
		}
	}

	return q, nil
}

// ApplyFilter applies only a filter to a goqu query
func ApplyFilter(
	q *goqu.SelectDataset,
	s *schema.Schema,
	filterStr string,
) (*goqu.SelectDataset, error) {
	return ApplyQuery(q, s, QueryParams{Filter: filterStr})
}

// ApplySort applies only a sort to a goqu query
func ApplySort(
	q *goqu.SelectDataset,
	s *schema.Schema,
	sortStr string,
) (*goqu.SelectDataset, error) {
	return ApplyQuery(q, s, QueryParams{Sort: sortStr})
}

// ValidateFilter validates a filter string without applying it
func ValidateFilter(s *schema.Schema, filterStr string) error {
	if filterStr == "" {
		return nil
	}

	pl := planner.New(s)
	p := parser.New(filterStr)
	expr, err := p.Parse()
	if err != nil {
		return fmt.Errorf("filter parse error: %w", err)
	}

	_, err = pl.Plan(expr)
	if err != nil {
		return fmt.Errorf("filter plan error: %w", err)
	}

	return nil
}

// ValidateSort validates a sort string without applying it
func ValidateSort(s *schema.Schema, sortStr string) error {
	if sortStr == "" {
		return nil
	}

	pl := planner.New(s)
	sortObj, err := sort.Parse(sortStr)
	if err != nil {
		return fmt.Errorf("sort parse error: %w", err)
	}

	_, err = pl.ResolveSort(sortObj.Orderings)
	if err != nil {
		return fmt.Errorf("sort resolution error: %w", err)
	}

	return nil
}

// ValidateQuery validates both filter and sort strings
func ValidateQuery(s *schema.Schema, params QueryParams) error {
	if err := ValidateFilter(s, params.Filter); err != nil {
		return err
	}
	if err := ValidateSort(s, params.Sort); err != nil {
		return err
	}
	return nil
}
