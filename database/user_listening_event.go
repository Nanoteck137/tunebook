package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
)

var createUserListeningEventId = createIdGenerator(32)

type UserListeningEvent struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	UserId  string `db:"user_id"`
	TrackId string `db:"track_id"`

	ListenedAt int64   `db:"listened_at"`
	Percent    float64 `db:"percent"`
	PositionMs int64   `db:"position_ms"`
	Source     string  `db:"source"`
}

func UserListeningEventQuery() *goqu.SelectDataset {
	query := dialect.From("user_listening_events").
		Select(
			"user_listening_events.rowid",

			"user_listening_events.id",

			"user_listening_events.user_id",
			"user_listening_events.track_id",

			"user_listening_events.listened_at",
			"user_listening_events.percent",
			"user_listening_events.position_ms",
			"user_listening_events.source",
		)

	return query
}

func (db DB) GetUserListeningEventById(ctx context.Context, id string) (UserListeningEvent, error) {
	query := UserListeningEventQuery().
		Where(
			goqu.I("user_listening_events.id").Eq(id),
		)

	return ember.Single[UserListeningEvent](db.db, ctx, query)
}

type CreateUserListeningEventParams struct {
	Id string

	UserId  string
	TrackId string

	ListenedAt int64
	Percent    float64
	PositionMs int64
	Source     string
}

func (db DB) CreateUserListeningEvent(ctx context.Context, params CreateUserListeningEventParams) (string, error) {
	if params.Id == "" {
		params.Id = createUserListeningEventId()
	}

	query := dialect.Insert("user_listening_events").Rows(goqu.Record{
		"id": params.Id,

		"user_id":  params.UserId,
		"track_id": params.TrackId,

		"listened_at": params.ListenedAt,
		"percent":     params.Percent,
		"position_ms": params.PositionMs,
		"source":      params.Source,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

func (db DB) DeleteUserListeningEvent(ctx context.Context, id string) error {
	query := dialect.Delete("user_listening_events").
		Where(
			goqu.I("user_listening_events.id").Eq(id),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
