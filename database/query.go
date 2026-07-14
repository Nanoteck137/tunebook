package database

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/parser"
	"github.com/nanoteck137/tunebook/tools/query/planner"
	"github.com/nanoteck137/tunebook/tools/query/schema"
	"github.com/nanoteck137/tunebook/tools/query/sort"
	querysql "github.com/nanoteck137/tunebook/tools/query/sql"
)

type FilterError struct {
	Op  string
	Err error
}

func (e *FilterError) Error() string {
	return fmt.Sprintf("filter %s error: %s", e.Op, e.Err)
}

func (e *FilterError) Unwrap() error { return e.Err }

type SortError struct {
	Op  string
	Err error
}

func (e *SortError) Error() string {
	return fmt.Sprintf("sort %s error: %s", e.Op, e.Err)
}

func (e *SortError) Unwrap() error { return e.Err }

type QueryError struct {
	Filter *FilterError
	Sort   *SortError
}

func (e *QueryError) Error() string {
	if e.Filter != nil && e.Sort != nil {
		return fmt.Sprintf("query errors: %s; %s", e.Filter.Error(), e.Sort.Error())
	}
	if e.Filter != nil {
		return e.Filter.Error()
	}
	if e.Sort != nil {
		return e.Sort.Error()
	}
	return "query error"
}

func ValidateFilter(schema *schema.Schema, filter string) error {
	if filter == "" {
		return nil
	}

	pl := planner.New(schema)

	p := parser.New(filter)
	expr, err := p.Parse()
	if err != nil {
		return &FilterError{Op: "parse", Err: err}
	}

	_, err = pl.Plan(expr)
	if err != nil {
		return &FilterError{Op: "plan", Err: err}
	}

	return nil
}

type QueryParams struct {
	Filter string
	Sort   string
}

func ApplyQuery(
	q *goqu.SelectDataset,
	s *schema.Schema,
	params QueryParams,
) (*goqu.SelectDataset, error) {
	pl := planner.New(s)
	compiler := querysql.NewCompiler()

	var filterErr *FilterError
	var sortErr *SortError

	// Validate filter
	var filterWhere *goqu.Expression
	if params.Filter != "" {
		p := parser.New(params.Filter)
		expr, err := p.Parse()
		if err != nil {
			filterErr = &FilterError{Op: "parse", Err: err}
		} else {
			plan, err := pl.Plan(expr)
			if err != nil {
				filterErr = &FilterError{Op: "plan", Err: err}
			} else {
				result, err := compiler.Compile(plan)
				if err != nil {
					filterErr = &FilterError{Op: "compile", Err: err}
				} else if result.Where != nil {
					filterWhere = &result.Where
				}
			}
		}
	}

	// Validate sort
	var sortOrder []exp.OrderedExpression
	if params.Sort != "" {
		sortObj, err := sort.Parse(params.Sort)
		if err != nil {
			sortErr = &SortError{Op: "parse", Err: err}
		} else {
			resolvedOrderings, err := pl.ResolveSort(sortObj.Orderings)
			if err != nil {
				sortErr = &SortError{Op: "resolution", Err: err}
			} else if len(resolvedOrderings) > 0 {
				plan := &query.Plan{OrderBy: resolvedOrderings}
				result, err := compiler.Compile(plan)
				if err != nil {
					sortErr = &SortError{Op: "compile", Err: err}
				} else {
					sortOrder = result.Order
				}
			}
		}
	}

	// Return errors if any
	if filterErr != nil || sortErr != nil {
		return nil, &QueryError{Filter: filterErr, Sort: sortErr}
	}

	// Apply filter
	if filterWhere != nil {
		q = q.Where(*filterWhere)
	}

	// Apply sort
	if len(sortOrder) > 0 {
		q = q.Order(sortOrder...)
	}

	return q, nil
}
