package usecase

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	"github.com/Eretic431/datingTelegramBot/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestUsecase_GetUserByIdOrNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	userId := "user_id"
	expectedUser := &models.User{Id: userId}

	ctx := context.Background()
	usersRepo.EXPECT().
		GetByUserId(ctx, userId).
		Return(expectedUser, nil).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	user, err := usecase.GetUserByIdOrNil(ctx, userId)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedUser, user)
}

func TestUsecase_GetUserByIdOrNil_ShouldReturnNilOnErrNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	userId := "user_id"

	ctx := context.Background()
	usersRepo.EXPECT().
		GetByUserId(ctx, userId).
		Return(nil, models.ErrNoRecord).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	user, err := usecase.GetUserByIdOrNil(ctx, userId)
	assert.Nil(t, err)
	assert.Nil(t, user)
}

func TestUsecase_GetUserByIdOrNil_ShouldReturnSameErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	userId := "user_id"
	expectedErr := errors.New("some err")

	ctx := context.Background()
	usersRepo.EXPECT().
		GetByUserId(ctx, userId).
		Return(nil, expectedErr).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	user, err := usecase.GetUserByIdOrNil(ctx, userId)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, expectedErr))
	assert.Nil(t, user)
}
