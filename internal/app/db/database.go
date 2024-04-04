//go:generate mockgen -source=./database.go -destination=./mocks/database.go -package=mocks

package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	cluster *pgxpool.Pool
}

func newDatabase(cluster *pgxpool.Pool) *PostgresDatabase {
	return &PostgresDatabase{cluster: cluster}
}

func (db PostgresDatabase) Close() {
	db.cluster.Close()
}

func (db PostgresDatabase) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.cluster, dest, query, args...)
}

func (db PostgresDatabase) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.cluster, dest, query, args...)
}

func (db PostgresDatabase) Exec(ctx context.Context, query string, args ...interface{}) (CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}

func (db PostgresDatabase) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}
