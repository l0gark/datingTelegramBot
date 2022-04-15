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

type UserRepository struct {
	DB PgxPoolIface
}

var _ internal.UsersRepository = &UserRepository{}

func NewUserRepository(DB PgxPoolIface) internal.UsersRepository {
	return &UserRepository{DB: DB}
}

func (ur *UserRepository) Add(ctx context.Context, user *models.User) error {
	tx, err := ur.DB.Begin(ctx)
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

	stage := user.Stage
	if stage == 0 {
		stage = -1
	}

	query := "INSERT INTO users (id, name, sex, age, description, city, image, started, stage, chat_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"

	if _, err = tx.Exec(ctx, query,
		user.Id,
		user.Name,
		user.Sex,
		user.Age,
		user.Description,
		user.City,
		user.Image,
		user.Started,
		stage,
		user.ChatId,
	); err != nil {
		pgErr := &pgconn.PgError{}

		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			return models.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (ur *UserRepository) GetByUserId(ctx context.Context, userId string) (user *models.User, err error) {
	tx, err := ur.DB.Begin(ctx)
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

	user = &models.User{}
	query := "SELECT id, name, sex, age, description, city, image, started, stage, chat_id FROM users WHERE id=$1;"

	if err := pgxscan.Get(ctx, tx,
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
	tx, err := ur.DB.Begin(ctx)
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

	query := "UPDATE users SET name=$2, sex=$3, age=$4, description=$5, city=$6, image=$7, started=$8, stage=$9, chat_id=$10 WHERE id=$1;"

	tag, err := tx.Exec(ctx, query,
		user.Id,
		user.Name,
		user.Sex,
		user.Age,
		user.Description,
		user.City,
		user.Image,
		user.Started,
		user.Stage,
		user.ChatId,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		err = models.ErrNoRecord
		return err
	}

	return nil
}

func (ur *UserRepository) DeleteByUserId(ctx context.Context, userId string) (err error) {
	tx, err := ur.DB.Begin(ctx)
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

	query := "DELETE FROM users WHERE id = $1;"
	tag, err := tx.Exec(ctx, query, userId)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return models.ErrNoRecord
	}

	return nil
}

func (ur *UserRepository) GetNextUser(ctx context.Context, userId string, sex bool) (*models.User, error) {
	tx, err := ur.DB.Begin(ctx)
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

	user := &models.User{}
	query := "SELECT id, name, sex, age, description, city, image, started, stage, chat_id FROM users" +
		" WHERE id IN (" +
		" SELECT user_ids.id as user_id FROM likes as likes2 " +
		" 	RIGHT JOIN ( " +
		"		SELECT users.id as id FROM users " +
		"			LEFT JOIN ( " +
		"				SELECT * FROM likes WHERE likes.from_id != $1" +
		"			) likes1 ON users.id = likes1.to_id" +
		"				 WHERE users.id != $1 " +
		"						AND users.id NOT IN (" +
		"							SELECT to_id as id FROM likes WHERE from_id = $1" +
		"						) " +
		"						AND users.sex != $2" +
		"	) user_ids ON likes2.from_id = user_ids.id AND likes2.to_id = $1 " +
		"	ORDER BY likes2.value DESC NULLS LAST" +
		"	LIMIT 1" +
		");"

	err = pgxscan.Get(ctx, tx,
		user,
		query,
		userId,
		sex,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return user, nil
}
