package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	userIdentitiesTbl = goqu.T("user_identities")
)

type UserIdentity struct {
	Provider   string `db:"provider"`
	ProviderId string `db:"provider_id"`

	UserId string `db:"user_id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func UserIdentityQuery() *goqu.SelectDataset {
	query := dialect.From(userIdentitiesTbl).
		Select(
			userIdentitiesTbl.Col("provider"),
			userIdentitiesTbl.Col("provider_id"),

			userIdentitiesTbl.Col("user_id"),

			userIdentitiesTbl.Col("created"),
			userIdentitiesTbl.Col("updated"),
		)

	return query
}

type CreateUserIdentityParams struct {
	Provider   string
	ProviderId string

	UserId string

	Created int64
	Updated int64
}

func (db DB) CreateUserIdentity(
	ctx context.Context,
	params CreateUserIdentityParams,
) error {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	query := dialect.Insert(userIdentitiesTbl).
		Rows(goqu.Record{
			"provider":    params.Provider,
			"provider_id": params.ProviderId,

			"user_id": params.UserId,

			"created": created,
			"updated": updated,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetUserIdentity(
	ctx context.Context,
	provider, providerId string,
) (UserIdentity, error) {
	query := UserIdentityQuery().
		Where(
			userIdentitiesTbl.Col("provider").Eq(provider),
			userIdentitiesTbl.Col("provider_id").Eq(providerId),
		)

	return Single[UserIdentity](db, ctx, query)
}
