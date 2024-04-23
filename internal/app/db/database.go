//go:generate mockgen -source=./database.go -destination=./mocks/database.go -package=mocks

package db

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type CommandTag interface {
	RowsAffected() int64
}

type Database interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type PostgresDatabase struct {
	provider QueryEngineProvider
}

func NewDatabase(provider QueryEngineProvider) *PostgresDatabase {
	return &PostgresDatabase{provider: provider}
}

func (db PostgresDatabase) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.provider.GetQueryEngine(ctx), dest, query, args...)
}

func (db PostgresDatabase) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.provider.GetQueryEngine(ctx), dest, query, args...)
}

func (db PostgresDatabase) Exec(ctx context.Context, query string, args ...interface{}) (CommandTag, error) {
	return db.provider.GetQueryEngine(ctx).Exec(ctx, query, args...)
}

func (db PostgresDatabase) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.provider.GetQueryEngine(ctx).QueryRow(ctx, query, args...)
}
