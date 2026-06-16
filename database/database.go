package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/nanoteck137/tunebook/database/migrations"

	"github.com/mattn/go-sqlite3"
)

var ErrItemNotFound = errors.New("database: item not found")
var ErrItemAlreadyExists = errors.New("database: item already exists")

type Query interface {
	ToSQL() (string, []any, error)
}

var dialect = goqu.Dialect("sqlite3")

type Executor interface {
	Query(ctx context.Context, query Query) (*sql.Rows, error)
	QueryRow(ctx context.Context, query Query) (*sql.Row, error)
	Exec(ctx context.Context, query Query) (sql.Result, error)

	Single(ctx context.Context, query Query, dest any) error
	Multiple(ctx context.Context, query Query, dest any) error
}

// Simple wrappers for retriving single row from the database
func Single[T any](
	executor Executor,
	ctx context.Context,
	query Query,
) (T, error) {
	var res T
	err := executor.Single(ctx, query, &res)
	return res, err
}

// Simple wrappers for retriving multiple row from the database
func Multiple[T any](
	executor Executor,
	ctx context.Context,
	query Query,
) ([]T, error) {
	var res []T
	err := executor.Multiple(ctx, query, &res)
	return res, err
}

type DB struct {
	Executor
}

var _ Executor = (*DBExecutor)(nil)

type DBExecutor struct {
	handle *sqlx.DB
}

func (db *DBExecutor) Exec(ctx context.Context, query Query) (sql.Result, error) {
	sql, params, err := query.ToSQL()
	if err != nil {
		return nil, handleErr(err)
	}

	res, err := db.handle.ExecContext(ctx, sql, params...)
	if err != nil {
		return nil, handleErr(err)
	}

	return res, nil
}

func (db *DBExecutor) Multiple(ctx context.Context, query Query, dest any) error {
	sql, params, err := query.ToSQL()
	if err != nil {
		return handleErr(err)
	}

	err = db.handle.SelectContext(ctx, dest, sql, params...)
	if err != nil {
		return handleErr(err)
	}

	return nil
}

func (db *DBExecutor) Query(ctx context.Context, query Query) (*sql.Rows, error) {
	sql, params, err := query.ToSQL()
	if err != nil {
		return nil, handleErr(err)
	}

	res, err := db.handle.QueryContext(ctx, sql, params...)
	if err != nil {
		return nil, handleErr(err)
	}

	return res, nil
}

func (db *DBExecutor) QueryRow(ctx context.Context, query Query) (*sql.Row, error) {
	sql, params, err := query.ToSQL()
	if err != nil {
		return nil, handleErr(err)
	}

	res := db.handle.QueryRowContext(ctx, sql, params...)
	return res, nil
}

func (db *DBExecutor) Single(ctx context.Context, query Query, dest any) error {
	sql, params, err := query.ToSQL()
	if err != nil {
		return handleErr(err)
	}

	err = db.handle.GetContext(ctx, dest, sql, params...)
	if err != nil {
		return handleErr(err)
	}

	return nil
}

type Database struct {
	DB

	handle *sqlx.DB
}

func (db *Database) RunMigrateUp() error {
	return migrations.RunMigrateUp(db.handle.DB)
}

func (db *Database) RunMigrateDown() error {
	return migrations.RunMigrateDown(db.handle.DB)
}

func (db *Database) Close() error {
	return db.Close()
}

func (db *Database) Begin() (Tx, error) {
	tx, err := db.handle.Beginx()
	if err != nil {
		return Tx{}, err
	}

	return Tx{
		DB: DB{
			Executor: &TxExecutor{
				Tx: tx,
			},
		},
	}, nil
}

func Open(dbFile string) (*Database, error) {
	dbUrl := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=ON&_serialized=1&_synchronous=NORMAL", dbFile)

	driver := "sqlite3"

	db, err := sql.Open(driver, dbUrl)
	if err != nil {
		return nil, err
	}

	handle := sqlx.NewDb(db, driver)

	return &Database{
		handle: handle,
		DB: DB{
			Executor: &DBExecutor{
				handle: handle,
			},
		},
	}, nil
}

type Tx struct {
	DB

	handle *sqlx.Tx
}

func (tx *Tx) Commit() error {
	return tx.handle.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.handle.Rollback()
}

var _ Executor = (*TxExecutor)(nil)

type TxExecutor struct {
	*sqlx.Tx
}

func (tx *TxExecutor) Exec(ctx context.Context, query Query) (sql.Result, error) {
	sql, params, err := query.ToSQL()
	if err != nil {
		return nil, handleErr(err)
	}

	res, err := tx.ExecContext(ctx, sql, params...)
	if err != nil {
		return nil, handleErr(err)
	}

	return res, nil
}

func (tx *TxExecutor) Multiple(ctx context.Context, query Query, dest any) error {
	sql, params, err := query.ToSQL()
	if err != nil {
		return handleErr(err)
	}

	err = tx.SelectContext(ctx, dest, sql, params...)
	if err != nil {
		return handleErr(err)
	}

	return nil
}

func (tx *TxExecutor) Query(ctx context.Context, query Query) (*sql.Rows, error) {
	sql, params, err := query.ToSQL()
	if err != nil {
		return nil, handleErr(err)
	}

	res, err := tx.QueryContext(ctx, sql, params...)
	if err != nil {
		return nil, handleErr(err)
	}

	return res, nil
}

func (tx *TxExecutor) QueryRow(ctx context.Context, query Query) (*sql.Row, error) {
	sql, params, err := query.ToSQL()
	if err != nil {
		return nil, handleErr(err)
	}

	res := tx.QueryRowContext(ctx, sql, params...)
	return res, nil
}

func (tx *TxExecutor) Single(ctx context.Context, query Query, dest any) error {
	sql, params, err := query.ToSQL()
	if err != nil {
		return handleErr(err)
	}

	err = tx.GetContext(ctx, dest, sql, params...)
	if err != nil {
		return handleErr(err)
	}

	return nil
}

func handleErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrItemNotFound
	}

	var e sqlite3.Error
	if errors.As(err, &e) {
		switch e.ExtendedCode {
		case sqlite3.ErrConstraintPrimaryKey:
			return ErrItemAlreadyExists
		}
	}

	return err
}
