package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/tunebook/tools/filter"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
	"github.com/nrednav/cuid2"
)

type RawQuery struct {
	Query  string
	Params []any
}

func (r RawQuery) ToSQL() (string, []any, error) {
	return r.Query, r.Params, nil
}

type Change[T any] struct {
	Value   T
	Changed bool
}

func addToRecord[T any](record goqu.Record, name string, change Change[T]) {
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

func SqlGroupConcat(col any, seperator string) exp.SQLFunctionExpression {
	return goqu.Func("group_concat", col, seperator)
}

func applyFilterParams(
	params types.QueryParams,
	adapter filter.ResolverAdapter,
	query *goqu.SelectDataset,
) (*goqu.SelectDataset, error) {
	resolver := filter.New(adapter)

	query, err := applyFilter(query, resolver, params.Filter)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, params.Sort)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func applyFilterParamsCustom(
	params types.QueryParams,
	adapter filter.ResolverAdapter,
	query *goqu.SelectDataset,
	customWhere exp.Expression,
) (*goqu.SelectDataset, error) {
	resolver := filter.New(adapter)

	query, err := applyFilterCustom(
		query, resolver, params.Filter, customWhere)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, params.Sort)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func buildPage(
	ctx context.Context,
	db DB,
	params types.PageParams,
	query *goqu.SelectDataset,
	countCol any,
) (types.Page, error) {
	countQuery := query.Select(goqu.COUNT(countCol))

	totalItems, err := Single[int](db, ctx, countQuery)
	if err != nil {
		return types.Page{}, err
	}

	return types.Page{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
		TotalPages: utils.TotalPages(params.PerPage, totalItems),
	}, nil
}

func applyPageParams(
	params types.PageParams,
	query *goqu.SelectDataset,
) *goqu.SelectDataset {
	return query.
		Limit(uint(params.PerPage)).
		Offset(uint(params.Page * params.PerPage))
}
