package postgres

import "github.com/jackc/pgx/v4/pgxpool"

type UserRepository struct {
	DB *pgxpool.Pool
}
