package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPsqlPool(c *Config) (PgxPoolIface, func(), error) {
	pool, err := pgxpool.Connect(context.Background(), c.PostgresUrl)
	if err != nil {
		return nil, nil, err
	}

	return pool, pool.Close, nil
}

var _ PgxPoolIface = &pgxpool.Pool{}

type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close()
}
