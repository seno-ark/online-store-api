package database

import (
	"context"
	"database/sql"
	"fmt"
	"online-store/pkg/config"

	_ "github.com/lib/pq"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func Postgres(conf *config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.Postgres.User,
		conf.Postgres.Pass,
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.DB,
	)

	return sql.Open("postgres", dataSourceName)
}
