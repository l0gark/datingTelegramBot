package postgres

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func Test1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//pool, err := pgxmock.NewPool()
	//if err != nil {
	//	// TODO
	//	return
	//}
	//
	//
	//users := NewUserRepository(pool)
	//
	//pool.ExpectExec("UPDATE ").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	//
	//
}
