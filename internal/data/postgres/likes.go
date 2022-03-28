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

func NewLikeRepository(DB PgxPoolIface) *LikeRepository {
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

func (lr *LikeRepository) Update(ctx context.Context, like *models.Like) error {
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

	query := "DELETE FROM likes WHERE id = $1"
	if _, err := tx.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (lr *LikeRepository) GetNewMatches(ctx context.Context, userId string) ([]models.User, error) {
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

	var users []models.User

	query := "SELECT users.id, users.name, users.sex, users.age, users.description, users.city, users.image FROM users" +
		" JOIN" +
		" (SELECT likes1.from_id FROM likes as likes1" +
		" 	JOIN likes as likes2 ON " +
		" 		likes1.from_id = likes2.to_id AND " +
		"		likes1.to_id = likes2.from_id" +
		" 	WHERE (likes1.value = true) AND (likes2.value = true) AND (likes1.to_id = $1)" +
		"	ORDER BY likes1.id) likes3 ON" +
		" users.id = likes3.from_id;"

	if err := pgxscan.Select(ctx, tx, users, query, userId); err != nil {
		return nil, err
	}

	return users, err
}
