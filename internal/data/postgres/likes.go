package postgres

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

type LikeRepository struct {
	DB PgxPoolIface
}

var _ internal.LikesRepository = &LikeRepository{}

func NewLikeRepository(DB PgxPoolIface) internal.LikesRepository {
	return &LikeRepository{DB: DB}
}

func (lr *LikeRepository) Add(ctx context.Context, like *models.Like) (err error) {
	tx, err := lr.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(ctx)
		default:
			_ = tx.Rollback(ctx)
		}
	}()

	query := "INSERT INTO likes (from_id, to_id, value) VALUES ($1, $2, $3);"

	if _, err := tx.Exec(ctx, query,
		like.FromId,
		like.ToId,
		like.Value,
	); err != nil {
		pgErr := &pgconn.PgError{}

		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			return models.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (lr *LikeRepository) Get(ctx context.Context, userFromId string, userToId string) (like *models.Like, err error) {
	tx, err := lr.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(ctx)
		default:
			_ = tx.Rollback(ctx)
		}
	}()

	like = &models.Like{}
	query := "SELECT id, from_id, to_id, value FROM likes WHERE from_id=$1 AND to_id=$2"

	if err := pgxscan.Get(ctx, tx, like, query, userFromId, userToId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return like, nil
}

func (lr *LikeRepository) Update(ctx context.Context, like *models.Like) (err error) {
	tx, err := lr.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(ctx)
		default:
			_ = tx.Rollback(ctx)
		}
	}()

	query := "UPDATE likes SET from_id=$2, to_id=$3, value=$4 WHERE id=$1"

	tag, err := tx.Exec(ctx, query,
		like.Id,
		like.FromId,
		like.ToId,
		like.Value,
	)

	if err != nil {
		pgErr := &pgconn.PgError{}
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

func (lr *LikeRepository) Delete(ctx context.Context, id int64) (err error) {
	tx, err := lr.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(ctx)
		default:
			_ = tx.Rollback(ctx)
		}
	}()

	query := "DELETE FROM likes WHERE id = $1;"
	tag, err := tx.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return models.ErrNoRecord
	}

	return nil
}

func (lr *LikeRepository) DeleteAll(ctx context.Context) (err error) {
	tx, err := lr.DB.Begin(ctx)
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(ctx)
		default:
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, "DELETE FROM likes;")

	return
}
