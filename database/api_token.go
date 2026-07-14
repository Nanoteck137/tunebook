package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	createApiTokenId = createIdGenerator(32)

	apiTokensTbl = goqu.T("api_tokens")
)

type ApiToken struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func ApiTokenQuery() *goqu.SelectDataset {
	query := dialect.From(apiTokensTbl).
		Select(
			apiTokensTbl.Col("id"),
			apiTokensTbl.Col("user_id"),

			apiTokensTbl.Col("name"),

			apiTokensTbl.Col("updated"),
			apiTokensTbl.Col("created"),
		)

	return query
}

func (db DB) GetApiTokenById(
	ctx context.Context,
	tokenId string,
) (ApiToken, error) {
	query := ApiTokenQuery().
		Where(apiTokensTbl.Col("id").Eq(tokenId))

	return Single[ApiToken](db, ctx, query)
}

func (db DB) GetAllApiTokensForUser(
	ctx context.Context,
	userId string,
) ([]ApiToken, error) {
	query := ApiTokenQuery().
		Where(apiTokensTbl.Col("user_id").Eq(userId))

	return Multiple[ApiToken](db, ctx, query)
}

type CreateApiTokenParams struct {
	Id     string
	UserId string
	Name   string

	Created int64
	Updated int64
}

func (db DB) CreateApiToken(
	ctx context.Context,
	params CreateApiTokenParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createApiTokenId()
	}

	query := dialect.Insert(apiTokensTbl).Rows(goqu.Record{
		"id":      params.Id,
		"user_id": params.UserId,

		"name": params.Name,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

func (db DB) DeleteApiToken(ctx context.Context, tokenId string) error {
	query := dialect.Delete(apiTokensTbl).
		Where(apiTokensTbl.Col("id").Eq(tokenId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
