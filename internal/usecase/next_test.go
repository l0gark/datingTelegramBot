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

func TestUsecase_HandleCommandNext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	var expectedChatId int64 = 1
	inputUser := &models.User{
		Id:  "123",
		Sex: false,
	}

	expectedUser := &models.User{
		Id:    "123",
		Sex:   false,
		Image: "123",
	}

	usersRepo.EXPECT().
		GetNextUser(gomock.Any(), inputUser.Id, inputUser.Sex).
		Return(expectedUser, nil).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	messageConfig, err := usecase.HandleCommandNext(context.Background(), expectedChatId, inputUser)
	assert.Nil(t, err)
	assert.NotNil(t, messageConfig)

	photoCfg, ok := messageConfig.(tgbotapi.PhotoConfig)
	assert.NotNil(t, photoCfg)
	assert.True(t, ok)
	if !ok {
		return
	}

	assert.EqualValues(t, expectedChatId, photoCfg.ChatID)
	assert.EqualValues(t, tgbotapi.ModeMarkdown, photoCfg.ParseMode)
}

func TestUsecase_HandleCommandNextOnErrorNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	var expectedChatId int64 = 1
	inputUser := &models.User{
		Id:  "123",
		Sex: false,
	}

	usersRepo.EXPECT().
		GetNextUser(gomock.Any(), inputUser.Id, inputUser.Sex).
		Return(nil, models.ErrNoRecord).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	chattable, err := usecase.HandleCommandNext(context.Background(), expectedChatId, inputUser)
	assert.Nil(t, err)
	assert.NotNil(t, chattable)

	messageCfg, ok := chattable.(tgbotapi.MessageConfig)
	assert.NotNil(t, messageCfg)
	assert.True(t, ok)
	if !ok {
		return
	}

	assert.EqualValues(t, expectedChatId, messageCfg.ChatID)
}

func TestUsecase_HandleCommandNextOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	var expectedChatId int64 = 1
	expectedError := errors.New("some error")
	inputUser := &models.User{
		Id:  "123",
		Sex: false,
	}

	usersRepo.EXPECT().
		GetNextUser(gomock.Any(), inputUser.Id, inputUser.Sex).
		Return(nil, expectedError).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	chattable, err := usecase.HandleCommandNext(context.Background(), expectedChatId, inputUser)
	assert.True(t, errors.Is(err, expectedError))
	assert.NotNil(t, chattable)

	messageCfg, ok := chattable.(tgbotapi.MessageConfig)
	assert.NotNil(t, messageCfg)
	assert.True(t, ok)
}

func TestUsecase_HandleCommandNextOnUserNilWithoutFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	var expectedChatId int64 = 1
	inputUser := &models.User{
		Id:  "123",
		Sex: false,
	}

	usersRepo.EXPECT().
		GetNextUser(gomock.Any(), inputUser.Id, inputUser.Sex).
		Return(nil, nil).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	chattable, err := usecase.HandleCommandNext(context.Background(), expectedChatId, inputUser)
	assert.Nil(t, err)
	assert.NotNil(t, chattable)

	messageCfg, ok := chattable.(tgbotapi.MessageConfig)
	assert.NotNil(t, messageCfg)
	assert.True(t, ok)
}
