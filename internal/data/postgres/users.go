package postgres

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

// Add inserts user. How to get context? Call ctx := context.Background() and pass it to function below.
func (ur *UserRepository) Add(ctx context.Context, user *models.User) error {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := "INSERT INTO users (id) VALUES ($1);"

	if _, err := conn.Exec(ctx,
		query,
		user.Id); err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			return models.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (ur *UserRepository) GetByUserId(ctx context.Context, userId string) (*models.User, error) {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := "SELECT id FROM users WHERE id=$1"

	user := &models.User{}
	if err := pgxscan.Get(ctx, conn, user, query, userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return user, nil
}
