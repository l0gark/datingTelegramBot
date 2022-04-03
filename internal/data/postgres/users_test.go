package postgres

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Add(t *testing.T) {
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

func TestUserRepository_Add_ShouldReturnErrorAlreadyExistsOnUniqueViolationFailure(t *testing.T) {
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

func TestUserRepository_Add_ShouldReturnSameErrorOnFailure(t *testing.T) {
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

func TestUserRepository_UpdateByUserId(t *testing.T) {
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
	pool.ExpectExec("UPDATE users ").WithArgs(
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
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	pool.ExpectCommit()

	users := NewUserRepository(pool)

	if err := users.UpdateByUserId(context.Background(), &user); err != nil {
		t.Errorf("error was not expected while updating user: %s", err.Error())
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_UpdateByUserId_ShouldReturnErrNoRecordOnEmptyRawsAffected(t *testing.T) {
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
	pool.ExpectExec("UPDATE users ").WithArgs(
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
	).WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	pool.ExpectRollback()

	users := NewUserRepository(pool)

	if err := users.UpdateByUserId(context.Background(), &user); err != nil {
		assert.EqualValues(t, models.ErrNoRecord, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_UpdateByUserId_ShouldReturnSameErrorOnFailure(t *testing.T) {
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
	pool.ExpectExec("UPDATE users ").WithArgs(
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
	).WillReturnError(someError)
	pool.ExpectRollback()

	users := NewUserRepository(pool)

	if err := users.UpdateByUserId(context.Background(), &user); err != nil {
		assert.EqualValues(t, someError, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_DeleteByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	pool.ExpectBegin()
	pool.ExpectExec("DELETE FROM users ").WithArgs(
		"1",
	).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	pool.ExpectCommit()

	users := NewUserRepository(pool)

	if err := users.DeleteByUserId(context.Background(), "1"); err != nil {
		t.Errorf("error was not expected while deleting user: %s", err.Error())
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_DeleteByUserId_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	someError := errors.New("some error")

	pool.ExpectBegin()
	pool.ExpectExec("DELETE FROM users ").WithArgs(
		"1",
	).WillReturnError(someError)
	pool.ExpectRollback()

	users := NewUserRepository(pool)

	if err := users.DeleteByUserId(context.Background(), "1"); err != nil {
		assert.EqualValues(t, someError, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_DeleteByUserId_ShouldReturnErrNoRecordIfRawsNoAffected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	pool.ExpectBegin()
	pool.ExpectExec("DELETE FROM users ").WithArgs(
		"1",
	).WillReturnResult(pgxmock.NewResult("DELETE", 0))
	pool.ExpectRollback()

	users := NewUserRepository(pool)

	if err := users.DeleteByUserId(context.Background(), "1"); err != nil {
		assert.EqualValues(t, models.ErrNoRecord, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetByUserIdShouldReturnErrNoRecordIfUserIsNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	pool.ExpectBegin()
	pool.ExpectQuery("^SELECT (.+) FROM users ").WithArgs(
		"1",
	).WillReturnError(pgx.ErrNoRows)
	pool.ExpectRollback()

	users := NewUserRepository(pool)

	if _, err := users.GetByUserId(context.Background(), "1"); err != nil {
		assert.EqualValues(t, models.ErrNoRecord, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetByUserIdShouldReturnRaws(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	user := &models.User{
		Id: "1",
	}

	pool.ExpectBegin()
	pool.ExpectQuery("^SELECT (.+) FROM users ").WithArgs(
		"1",
	).WillReturnRows(pgxmock.NewRows([]string{"id", "name", "sex", "age", "description", "city", "image", "started", "stage", "chat_id"}).AddRow(
		user.Id, user.Name, user.Sex, user.Age, user.Description, user.City, user.Image, user.Started, user.Stage, user.ChatId,
	))
	pool.ExpectCommit()

	users := NewUserRepository(pool)

	actualUser, err := users.GetByUserId(context.Background(), "1")
	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err.Error())
	}

	assert.EqualValues(t, user, actualUser)

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetNextUserShouldReturnErrNoRecordNoRaws(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	pool.ExpectBegin()
	pool.ExpectQuery("^SELECT (.+) FROM users ").WithArgs(
		"1",
		true,
	).WillReturnError(pgx.ErrNoRows)
	pool.ExpectRollback()

	users := NewUserRepository(pool)

	if _, err := users.GetNextUser(context.Background(), "1", true); err != nil {
		assert.EqualValues(t, models.ErrNoRecord, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetNextUserShouldReturnRaws(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	user := &models.User{
		Id:  "1",
		Sex: false,
	}

	pool.ExpectBegin()
	pool.ExpectQuery("^SELECT (.+) FROM users ").WithArgs(
		"1",
		true,
	).WillReturnRows(pgxmock.NewRows([]string{"id", "name", "sex", "age", "description", "city", "image", "started", "stage", "chat_id"}).AddRow(
		"2", user.Name, !user.Sex, user.Age, user.Description, user.City, user.Image, user.Started, user.Stage, user.ChatId,
	))
	pool.ExpectCommit()

	users := NewUserRepository(pool)

	actualUser, err := users.GetNextUser(context.Background(), "1", true)
	if err != nil {
		t.Errorf("error was not expected while updating user: %s", err.Error())
	}

	assert.NotEqualValues(t, user.Sex, actualUser.Sex)
	assert.NotEqualValues(t, user.Id, actualUser.Id)

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
