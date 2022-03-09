package postgres

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPsqlPool(c *Config) (*pgxpool.Pool, func(), error) {
	pool, err := pgxpool.Connect(context.Background(), c.PostgresUrl)
	if err != nil {
		return nil, nil, err
	}

	return pool, pool.Close, nil
}

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) (err error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
	Deallocate(ctx context.Context, name string) error
}

type PgxPoolIface interface {
	PgxIface
	pgx.Tx
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Close()
}
