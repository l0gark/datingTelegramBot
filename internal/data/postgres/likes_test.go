package postgres

import (
	"context"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock"
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
		t.Errorf("error was not expected while inserting user: %s", err.Error())
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
