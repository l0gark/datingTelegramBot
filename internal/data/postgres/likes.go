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

type LikeRepository struct {
	DB *pgxpool.Pool
}

func (lr *LikeRepository) Add(ctx context.Context, like *models.Like) error {
	conn, err := lr.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := "INSERT INTO likes (from_id, to_id, showed) VALUES ($1, $2, $3);"

	if _, err := conn.Exec(ctx, query, like.FromId, like.ToId, like.Showed); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			return models.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (lr *LikeRepository) Get(ctx context.Context, id int64) (*models.Like, error) {
	conn, err := lr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	like := &models.Like{}
	query := "SELECT id, from_id, to_id, showed FROM likes WHERE id=$1"

	if err := pgxscan.Get(ctx, conn, like, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return like, nil
}

func (lr *LikeRepository) Update(ctx context.Context, like *models.Like) error {
	conn, err := lr.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	stmt := "UPDATE likes SET from_id=$2, to_id=$3, showed=$4 WHERE id=$1"

	tag, err := conn.Exec(ctx, stmt, like.Id, like.FromId, like.ToId, like.Showed)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			return models.ErrAlreadyExists
		}
		return err
	}

	if tag.RowsAffected() == 0 {
		return models.ErrNoRecord
	}

	return nil
}

func (lr *LikeRepository) Delete(ctx context.Context, id int64) error {
	conn, err := lr.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	stmt := "DELETE FROM likes WHERE id = $1"
	if _, err := conn.Exec(ctx, stmt, id); err != nil {
		return err
	}

	return nil
}
