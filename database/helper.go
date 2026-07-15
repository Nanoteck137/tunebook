package database

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

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

type JsonColumn[T any] struct {
	Data  T
	Valid bool
}

func (j *JsonColumn[T]) Scan(src any) error {
	var res T

	if src == nil {
		j.Data = res
		j.Valid = false
		return nil
	}

	switch value := src.(type) {
	case string:
		err := json.Unmarshal([]byte(value), &j.Data)
		if err != nil {
			return fmt.Errorf("jsoncolumn: failed to unmarshal data: %w", err)
		}

		j.Valid = true
	case []byte:
		err := json.Unmarshal(value, &j.Data)
		if err != nil {
			return fmt.Errorf("jsoncolumn: failed to unmarshal data: %w", err)
		}

		j.Valid = true
	default:
		return fmt.Errorf("jsoncolumn: unsupported src type %T", src)
	}

	return nil
}

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	raw, err := json.Marshal(j.Data)
	return raw, err
}

type KVStore map[string]string

func (kv KVStore) Serialize() (string, error) {
	b, err := json.Marshal(kv)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func DeserializeKVStore(data string) (KVStore, error) {
	kv := make(KVStore)
	if data == "" {
		return kv, nil
	}

	err := json.Unmarshal([]byte(data), &kv)
	if err != nil {
		return nil, err
	}

	return kv, nil
}

func (kv KVStore) Value() (driver.Value, error) {
	return kv.Serialize()
}

func (kv *KVStore) Scan(src any) error {
	if src == nil {
		*kv = make(KVStore)
		return nil
	}

	switch value := src.(type) {
	case string:
		r, err := DeserializeKVStore(value)
		if err != nil {
			return fmt.Errorf("kvstore: failed to deserialize store: %w", err)
		}

		*kv = r
	case []byte:
		r, err := DeserializeKVStore(string(value))
		if err != nil {
			return fmt.Errorf("kvstore: failed to deserialize store: %w", err)
		}

		*kv = r
	default:
		return fmt.Errorf("kvstore: unsupported src type %T", src)
	}

	return nil
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
