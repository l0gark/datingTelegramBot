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

func TestLikesRepository_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	like := &models.Like{
		FromId: "id1",
		ToId:   "id2",
		Value:  true,
	}

	pool.ExpectBegin()
	pool.ExpectExec("INSERT INTO likes ").WithArgs(
		like.FromId,
		like.ToId,
		like.Value,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	pool.ExpectCommit()

	likes := NewLikeRepository(pool)

	if err := likes.Add(context.Background(), like); err != nil {
		t.Errorf("error was not expected while inserting like: %s", err.Error())
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLikesRepository_Add_OnUniqueViolationReturnErrAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	like := &models.Like{
		FromId: "id1",
		ToId:   "id2",
		Value:  true,
	}

	pool.ExpectBegin()
	pool.ExpectExec("INSERT INTO likes ").WithArgs(
		like.FromId,
		like.ToId,
		like.Value,
	).WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
	pool.ExpectRollback()

	likes := NewLikeRepository(pool)

	if err := likes.Add(context.Background(), like); err != nil {
		assert.EqualValues(t, models.ErrAlreadyExists, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLikesRepository_Add_OnErrorShouldRollbackTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	like := &models.Like{
		FromId: "id1",
		ToId:   "id2",
		Value:  true,
	}

	someError := errors.New("some error")

	pool.ExpectBegin()
	pool.ExpectExec("INSERT INTO likes ").WithArgs(
		like.FromId,
		like.ToId,
		like.Value,
	).WillReturnError(someError)
	pool.ExpectRollback()

	likes := NewLikeRepository(pool)

	if err := likes.Add(context.Background(), like); err != nil {
		assert.EqualValues(t, someError, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLikeRepository_Get_ShouldReturnRaws(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	like := &models.Like{
		Id:     1,
		FromId: "id1",
		ToId:   "id2",
		Value:  true,
	}

	pool.ExpectBegin()
	pool.ExpectQuery("^SELECT (.+) FROM likes ").WithArgs(
		like.FromId, like.ToId,
	).WillReturnRows(pgxmock.NewRows([]string{"id", "from_id", "to_id", "value"}).AddRow(
		like.Id, like.FromId, like.ToId, like.Value,
	))
	pool.ExpectCommit()

	likes := NewLikeRepository(pool)

	actualLike, err := likes.Get(context.Background(), like.FromId, like.ToId)
	if err != nil {
		t.Errorf("error was not expected while getting like: %s", err.Error())
	}

	assert.EqualValues(t, like, actualLike)

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLikeRepository_Get_ShouldReturnErrNoRecordIfUserIsNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	like := &models.Like{
		FromId: "id1",
		ToId:   "id2",
	}

	pool.ExpectBegin()
	pool.ExpectQuery("^SELECT (.+) FROM likes ").WithArgs(
		like.FromId, like.ToId,
	).WillReturnError(pgx.ErrNoRows)
	pool.ExpectRollback()

	likes := NewLikeRepository(pool)

	if _, err := likes.Get(context.Background(), like.FromId, like.ToId); err != nil {
		assert.EqualValues(t, models.ErrNoRecord, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	like := &models.Like{
		Id:     1,
		FromId: "id1",
		ToId:   "id2",
		Value:  true,
	}

	pool.ExpectBegin()
	pool.ExpectExec("UPDATE likes ").WithArgs(
		like.Id,
		like.FromId,
		like.ToId,
		like.Value,
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	pool.ExpectCommit()

	likes := NewLikeRepository(pool)

	if err := likes.Update(context.Background(), like); err != nil {
		t.Errorf("error was not expected while updating like: %s", err.Error())
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Update_ShouldReturnErrNoRecordOnEmptyRawsAffected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	like := &models.Like{
		Id:     1,
		FromId: "id1",
		ToId:   "id2",
		Value:  true,
	}

	pool.ExpectBegin()
	pool.ExpectExec("UPDATE likes ").WithArgs(
		like.Id,
		like.FromId,
		like.ToId,
		like.Value,
	).WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	pool.ExpectRollback()

	likes := NewLikeRepository(pool)

	if err := likes.Update(context.Background(), like); err != nil {
		assert.EqualValues(t, models.ErrNoRecord, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Update_ShouldReturnSameErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("error was not expected while creating pool: %s", err.Error())
		return
	}
	defer pool.Close()

	like := &models.Like{
		Id:     1,
		FromId: "id1",
		ToId:   "id2",
		Value:  true,
	}

	someError := errors.New("some error")

	pool.ExpectBegin()
	pool.ExpectExec("UPDATE likes ").WithArgs(
		like.Id,
		like.FromId,
		like.ToId,
		like.Value,
	).WillReturnError(someError)
	pool.ExpectRollback()

	likes := NewLikeRepository(pool)

	if err := likes.Update(context.Background(), like); err != nil {
		assert.EqualValues(t, someError, err)
	} else {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
