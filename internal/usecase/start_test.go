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

func TestUsecase_HandleStart_IfStartedTrue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputMsg := &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 1}}

	usecase := NewUsecase(
		nil,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	msg, err := usecase.HandleStart(context.Background(), inputMsg, true)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	assert.EqualValues(t, inputMsg.Chat.ID, msg.ChatID)
	assert.EqualValues(t, tgbotapi.ModeMarkdown, msg.ParseMode)
}

func TestUsecase_HandleStart_IfStartedFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputMsg := &tgbotapi.Message{
		From: &tgbotapi.User{UserName: "username"},
		Chat: &tgbotapi.Chat{ID: 1},
	}

	usersRepo := mock.NewMockUsersRepository(ctrl)

	user := &models.User{
		Id:      inputMsg.From.UserName,
		Sex:     false,
		Started: true,
		Stage:   ProfileStageNone,
		ChatId:  inputMsg.Chat.ID,
	}

	usersRepo.EXPECT().
		UpdateByUserId(gomock.Any(), user).
		Return(nil).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		&tgbotapi.BotAPI{Self: tgbotapi.User{UserName: "botName"}},
		zaptest.NewLogger(t).Sugar(),
	)

	msg, err := usecase.HandleStart(context.Background(), inputMsg, false)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	assert.EqualValues(t, inputMsg.Chat.ID, msg.ChatID)
	assert.EqualValues(t, tgbotapi.ModeMarkdown, msg.ParseMode)
}

func TestUsecase_HandleStart_IfStartedFalse_ShouldReturnErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputMsg := &tgbotapi.Message{
		From: &tgbotapi.User{UserName: "username"},
		Chat: &tgbotapi.Chat{ID: 1},
	}

	usersRepo := mock.NewMockUsersRepository(ctrl)

	user := &models.User{
		Id:      inputMsg.From.UserName,
		Sex:     false,
		Started: true,
		Stage:   ProfileStageNone,
		ChatId:  inputMsg.Chat.ID,
	}

	expectedError := errors.New("some error")
	usersRepo.EXPECT().
		UpdateByUserId(gomock.Any(), user).
		Return(expectedError).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		&tgbotapi.BotAPI{Self: tgbotapi.User{UserName: "botName"}},
		zaptest.NewLogger(t).Sugar(),
	)

	msg, err := usecase.HandleStart(context.Background(), inputMsg, false)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, expectedError))
	assert.NotNil(t, msg)
	assert.EqualValues(t, tgbotapi.MessageConfig{}, msg)
}

func TestUsecase_IsStarted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputMsg := &tgbotapi.Message{
		From: &tgbotapi.User{UserName: "username"},
		Chat: &tgbotapi.Chat{ID: 1},
	}

	usersRepo := mock.NewMockUsersRepository(ctrl)

	expectedUser := &models.User{
		Started: true,
	}

	usersRepo.EXPECT().
		GetByUserId(gomock.Any(), inputMsg.From.UserName).
		Return(expectedUser, nil).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	started, err := usecase.IsStarted(context.Background(), inputMsg)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedUser.Started, started)
}

func TestUsecase_IsStarted_ShouldReturnSameErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputMsg := &tgbotapi.Message{
		From: &tgbotapi.User{UserName: "username"},
		Chat: &tgbotapi.Chat{ID: 1},
	}

	usersRepo := mock.NewMockUsersRepository(ctrl)

	expectedError := errors.New("some error")
	usersRepo.EXPECT().
		GetByUserId(gomock.Any(), inputMsg.From.UserName).
		Return(nil, expectedError).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	started, err := usecase.IsStarted(context.Background(), inputMsg)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, expectedError))
	assert.False(t, started)
}

func TestUsecase_IsStarted_ShouldAddUserIfNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputMsg := &tgbotapi.Message{
		From: &tgbotapi.User{UserName: "username"},
		Chat: &tgbotapi.Chat{ID: 1},
	}

	usersRepo := mock.NewMockUsersRepository(ctrl)

	user := &models.User{
		Id:      inputMsg.From.UserName,
		Sex:     false,
		Started: false,
		Stage:   ProfileStageNone,
		ChatId:  inputMsg.Chat.ID,
	}

	usersRepo.EXPECT().
		Add(gomock.Any(), user).
		After(
			usersRepo.EXPECT().
				GetByUserId(gomock.Any(), inputMsg.From.UserName).
				Return(nil, models.ErrNoRecord).
				Times(1),
		).
		Return(nil).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	started, err := usecase.IsStarted(context.Background(), inputMsg)
	assert.Nil(t, err)
	assert.False(t, started)
}

func TestUsecase_IsStarted_ShouldAddUserIfNotExists_ShouldReturnSameErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputMsg := &tgbotapi.Message{
		From: &tgbotapi.User{UserName: "username"},
		Chat: &tgbotapi.Chat{ID: 1},
	}

	usersRepo := mock.NewMockUsersRepository(ctrl)

	user := &models.User{
		Id:      inputMsg.From.UserName,
		Sex:     false,
		Started: false,
		Stage:   ProfileStageNone,
		ChatId:  inputMsg.Chat.ID,
	}

	expectedError := errors.New("some error")

	usersRepo.EXPECT().
		Add(gomock.Any(), user).
		After(
			usersRepo.EXPECT().
				GetByUserId(gomock.Any(), inputMsg.From.UserName).
				Return(nil, models.ErrNoRecord).
				Times(1),
		).
		Return(expectedError).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	started, err := usecase.IsStarted(context.Background(), inputMsg)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, expectedError))
	assert.False(t, started)
}
