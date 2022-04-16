package usecase

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestUsecase_DeleteAll_ShouldReturnErrorOnUsersRepoFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)
	likesRepo := mock.NewMockLikesRepository(ctrl)

	ctx := context.Background()

	expectedErr := errors.New("some err")
	likesRepo.EXPECT().DeleteAll(ctx).Return(nil)
	usersRepo.EXPECT().DeleteAll(ctx).Return(expectedErr)

	usecase := NewUsecase(
		usersRepo,
		likesRepo,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	err := usecase.DeleteAll(ctx)
	assert.NotNil(t, err)
	assert.EqualValues(t, expectedErr, err)
}

func TestUsecase_DeleteAll_ShouldReturnErrorOnLikesRepoFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)
	likesRepo := mock.NewMockLikesRepository(ctrl)

	ctx := context.Background()

	expectedErr := errors.New("some err")
	likesRepo.EXPECT().DeleteAll(ctx).Return(expectedErr)

	usecase := NewUsecase(
		usersRepo,
		likesRepo,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	err := usecase.DeleteAll(ctx)
	assert.NotNil(t, err)
	assert.EqualValues(t, expectedErr, err)
}

func TestUsecase_AddTestUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	ctx := context.Background()

	usersRepo.EXPECT().Add(ctx, gomock.Any()).Return(nil)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	err := usecase.AddTestUser(ctx, false)
	assert.Nil(t, err)
}

func TestUsecase_AddTestUser_ShouldReturnErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	ctx := context.Background()

	expectedErr := errors.New("some err")
	usersRepo.EXPECT().Add(ctx, gomock.Any()).Return(expectedErr)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	err := usecase.AddTestUser(ctx, false)
	assert.NotNil(t, err)
	assert.EqualValues(t, expectedErr, err)
}

func TestUsecase_AddTestUserWithLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)
	likesRepo := mock.NewMockLikesRepository(ctrl)

	ctx := context.Background()

	usersRepo.EXPECT().Add(ctx, gomock.Any()).Return(nil)
	likesRepo.EXPECT().Add(ctx, gomock.Any()).Return(nil)

	usecase := NewUsecase(
		usersRepo,
		likesRepo,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	err := usecase.AddTestUserWithLike(ctx, false, "toId")
	assert.Nil(t, err)
}

func TestUsecase_AddTestUserWithLike_ShouldReturnErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)
	likesRepo := mock.NewMockLikesRepository(ctrl)

	ctx := context.Background()
	expectedErr := errors.New("some err")

	usersRepo.EXPECT().Add(ctx, gomock.Any()).Return(nil)
	likesRepo.EXPECT().Add(ctx, gomock.Any()).Return(expectedErr)

	usecase := NewUsecase(
		usersRepo,
		likesRepo,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	err := usecase.AddTestUserWithLike(ctx, false, "toId")
	assert.NotNil(t, err)
	assert.EqualValues(t, expectedErr, err)
}
