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

func (ur *UserRepository) Add(ctx context.Context, user *models.User) error {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := "INSERT INTO users (id, name, sex, age, description, city, image, started) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"

	if _, err := conn.Exec(ctx, query,
		user.Id,
		user.Name,
		user.Sex,
		user.Age,
		user.Description,
		user.City,
		user.Image,
		user.Started,
	); err != nil {
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

	user := &models.User{}
	query := "SELECT id, name, sex, age, description, city, image, started FROM users WHERE id=$1"

	if err := pgxscan.Get(ctx, conn,
		user,
		query,
		userId,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) UpdateByUserId(ctx context.Context, user *models.User) error {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := "UPDATE users SET name=$2, sex=$3, age=$4, description=$5, city=$6, image=$7, started=$8 WHERE id=$1"

	tag, err := conn.Exec(ctx, query,
		user.Id,
		user.Name,
		user.Sex,
		user.Age,
		user.Description,
		user.City,
		user.Image,
		user.Started,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return models.ErrNoRecord
	}

	return nil
}

func (ur *UserRepository) DeleteByUserId(ctx context.Context, userId string) error {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := "DELETE FROM users WHERE id = $1"
	if _, err := conn.Exec(ctx, query, userId); err != nil {
		return err
	}

	return nil
}
