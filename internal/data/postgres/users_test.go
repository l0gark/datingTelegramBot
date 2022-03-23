package postgres

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldInsertUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	user := models.User{}

	pool.ExpectBegin()
	pool.ExpectExec("INSERT INTO users ").WithArgs(
		user.Id,
		user.Name,
		user.Sex,
		user.Age,
		user.Description,
		user.City,
		user.Image,
		user.Started,
		user.ChatId,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	pool.ExpectCommit()

	users := NewUserRepository(pool)

	if err := users.Add(context.Background(), &user); err != nil {
		t.Errorf("error was not expected while inserting user: %s", err.Error())
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertUserShouldReturnErrorAlreadyExistsOnUniqueViolationFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	user := models.User{}

	pool.ExpectBegin()
	pool.ExpectExec("INSERT INTO users ").WithArgs(
		user.Id,
		user.Name,
		user.Sex,
		user.Age,
		user.Description,
		user.City,
		user.Image,
		user.Started,
		user.ChatId,
	).WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
	pool.ExpectRollback()

	users := NewUserRepository(pool)

	if err := users.Add(context.Background(), &user); err != nil {
		assert.EqualValues(t, models.ErrAlreadyExists, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertUserShouldReturnSameErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	user := models.User{}

	someError := errors.New("some error")

	pool.ExpectBegin()
	pool.ExpectExec("INSERT INTO users ").WithArgs(
		user.Id,
		user.Name,
		user.Sex,
		user.Age,
		user.Description,
		user.City,
		user.Image,
		user.Started,
		user.ChatId,
	).WillReturnError(someError)
	pool.ExpectRollback()

	users := NewUserRepository(pool)

	if err := users.Add(context.Background(), &user); err != nil {
		assert.EqualValues(t, someError, err)
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
