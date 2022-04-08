package usecase

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	"github.com/Eretic431/datingTelegramBot/internal/mock"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestUsecase_HasLikeWithTrueValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	likesRepo := mock.NewMockLikesRepository(ctrl)
	fromId := "user1"
	toId := "user2"
	ctx := context.Background()

	expectedLike := &models.Like{Value: true}

	likesRepo.EXPECT().
		Get(ctx, fromId, toId).
		Return(expectedLike, nil).
		Times(1)

	usecase := NewUsecase(
		nil,
		likesRepo,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	hasLike, err := usecase.HasLikeWithTrueValue(ctx, fromId, toId)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedLike.Value, hasLike)
}

func TestUsecase_HasLikeWithTrueValue_ShouldReturnOnErrNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	likesRepo := mock.NewMockLikesRepository(ctrl)
	fromId := "user1"
	toId := "user2"
	ctx := context.Background()

	likesRepo.EXPECT().
		Get(ctx, fromId, toId).
		Return(nil, models.ErrNoRecord).
		Times(1)

	usecase := NewUsecase(
		nil,
		likesRepo,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	hasLike, err := usecase.HasLikeWithTrueValue(ctx, fromId, toId)
	assert.Nil(t, err)
	assert.False(t, hasLike)
}

func TestUsecase_HasLikeWithTrueValue_ShouldReturnSameErrOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	likesRepo := mock.NewMockLikesRepository(ctrl)
	fromId := "user1"
	toId := "user2"
	ctx := context.Background()

	expectedErr := errors.New("some err")

	likesRepo.EXPECT().
		Get(ctx, fromId, toId).
		Return(nil, expectedErr).
		Times(1)

	usecase := NewUsecase(
		nil,
		likesRepo,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	hasLike, err := usecase.HasLikeWithTrueValue(ctx, fromId, toId)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, expectedErr))
	assert.False(t, hasLike)
}

func TestUsecase_CreateMatchMessages_ShouldReturnErrOnNilUsers(t *testing.T) {
	usecase := NewUsecase(
		nil,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	msg1, msg2, err := usecase.CreateMatchMessages(nil, nil)
	assert.Nil(t, msg1)
	assert.Nil(t, msg2)
	assert.NotNil(t, err)

	msg1, msg2, err = usecase.CreateMatchMessages(&models.User{}, nil)
	assert.Nil(t, msg1)
	assert.Nil(t, msg2)
	assert.NotNil(t, err)

	msg1, msg2, err = usecase.CreateMatchMessages(nil, &models.User{})
	assert.Nil(t, msg1)
	assert.Nil(t, msg2)
	assert.NotNil(t, err)
}

func TestUsecase_CreateMatchMessages(t *testing.T) {
	usecase := NewUsecase(
		nil,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	user1 := &models.User{
		Id:          "1",
		Name:        "name1",
		Age:         1,
		Description: "desc1",
		City:        "city1",
		Image:       "image1",
		ChatId:      1,
	}

	user2 := &models.User{
		Id:          "2",
		Name:        "name2",
		Age:         2,
		Description: "desc2",
		City:        "city2",
		Image:       "image2",
		ChatId:      2,
	}

	msg1, msg2, err := usecase.CreateMatchMessages(user1, user2)
	assert.Nil(t, err)
	assert.NotNil(t, msg1)
	assert.NotNil(t, msg2)
	photo1, ok := msg1.(tgbotapi.PhotoConfig)
	assert.True(t, ok)
	assert.NotNil(t, photo1)
	photo2, ok := msg2.(tgbotapi.PhotoConfig)
	assert.True(t, ok)
	assert.NotNil(t, photo2)

	assert.EqualValues(t, tgbotapi.ModeMarkdown, photo1.ParseMode)
	assert.EqualValues(t, tgbotapi.ModeMarkdown, photo2.ParseMode)

	assert.EqualValues(t, user1.ChatId, photo1.ChatID)
	assert.EqualValues(t, user2.ChatId, photo2.ChatID)
}
