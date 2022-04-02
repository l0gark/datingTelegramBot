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

func TestUsecase_HandleProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	inputMsg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 1},
	}

	user := &models.User{Name: "name"}

	usersRepo.EXPECT().
		UpdateByUserId(gomock.Any(), user).
		Return(nil).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	messageCfg, err := usecase.HandleProfile(context.Background(), inputMsg, user)
	assert.Nil(t, err)
	assert.NotNil(t, messageCfg)
	assert.EqualValues(t, inputMsg.Chat.ID, messageCfg.ChatID)
}

func TestUsecase_HandleProfile_ShouldReturnErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock.NewMockUsersRepository(ctrl)

	inputMsg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 1},
	}

	user := &models.User{Name: "name"}
	expectedError := errors.New("some error")

	usersRepo.EXPECT().
		UpdateByUserId(gomock.Any(), user).
		Return(expectedError).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	messageCfg, err := usecase.HandleProfile(context.Background(), inputMsg, user)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, expectedError))
	assert.NotNil(t, messageCfg)
	assert.EqualValues(t, tgbotapi.MessageConfig{}, messageCfg)
}

func TestUsecase_HandleFillingProfile_Correct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var chatId int64 = 1
	photoId := "NoImageData"

	usersRepo := mock.NewMockUsersRepository(ctrl)
	usersRepo.EXPECT().
		UpdateByUserId(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(MaxProfileStage)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	data := []string{"name", "1", "city", "description", "image"}

	for stage := 0; stage < MaxProfileStage; stage++ {
		inputText := data[stage]
		user := &models.User{Id: "id", Stage: stage}

		chattable, err := usecase.HandleFillingProfile(context.Background(), inputText, chatId, photoId, user)
		assert.Nil(t, err)
		assert.NotNil(t, chattable)
		msgCfg, ok := chattable.(tgbotapi.MessageConfig)
		assert.True(t, ok)
		assert.NotNil(t, msgCfg)
		assert.EqualValues(t, chatId, msgCfg.ChatID)
	}
}

func TestUsecase_HandleFillingProfile_Incorrect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var chatId int64 = 1
	photoId := "-"

	usersRepo := mock.NewMockUsersRepository(ctrl)
	//usersRepo.EXPECT().
	//	UpdateByUserId(gomock.Any(), gomock.Any()).
	//	Return(nil).
	//	Times(MaxProfileStage)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	data := []string{"", "", "", "", ""}

	for stage := 0; stage < MaxProfileStage; stage++ {
		inputText := data[stage]
		user := &models.User{Id: "id", Stage: stage}

		chattable, err := usecase.HandleFillingProfile(context.Background(), inputText, chatId, photoId, user)
		assert.Nil(t, err)
		assert.NotNil(t, chattable)
		msgCfg, ok := chattable.(tgbotapi.MessageConfig)
		assert.True(t, ok)
		assert.NotNil(t, msgCfg)
		assert.EqualValues(t, chatId, msgCfg.ChatID)
	}
}

func TestUsecase_HandleFillingProfile_StageNameCorrect_ShouldReturnSameErrorOnUpdateFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputText := "text"
	var chatId int64 = 1
	photoId := "photoId"
	user := &models.User{Id: "id", Stage: 0}

	usersRepo := mock.NewMockUsersRepository(ctrl)
	expectedError := errors.New("some error")
	usersRepo.EXPECT().
		UpdateByUserId(gomock.Any(), user).
		Return(expectedError).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	chattable, err := usecase.HandleFillingProfile(context.Background(), inputText, chatId, photoId, user)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, expectedError))
	assert.NotNil(t, chattable)
	msgCfg, ok := chattable.(tgbotapi.MessageConfig)
	assert.True(t, ok)
	assert.NotNil(t, msgCfg)
	assert.EqualValues(t, tgbotapi.MessageConfig{}, msgCfg)
}

func TestUsecase_HandleFillingProfile_StageNameIncorrect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputText := ""
	var chatId int64 = 1
	photoId := "photoId"
	user := &models.User{Id: "id", Stage: 0}

	usersRepo := mock.NewMockUsersRepository(ctrl)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	chattable, err := usecase.HandleFillingProfile(context.Background(), inputText, chatId, photoId, user)
	assert.Nil(t, err)
	assert.NotNil(t, chattable)
	msgCfg, ok := chattable.(tgbotapi.MessageConfig)
	assert.True(t, ok)
	assert.NotNil(t, msgCfg)
	assert.EqualValues(t, chatId, msgCfg.ChatID)

}

func TestUsecase_HandleFillingProfile_StageSex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputText := "М"
	var chatId int64 = 1
	photoId := "photoId"
	user := &models.User{Id: "id", Stage: MaxProfileStage}

	usersRepo := mock.NewMockUsersRepository(ctrl)

	usersRepo.EXPECT().
		UpdateByUserId(gomock.Any(), user).
		Return(nil).
		Times(1)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	chattable, err := usecase.HandleFillingProfile(context.Background(), inputText, chatId, photoId, user)
	assert.Nil(t, err)
	assert.NotNil(t, chattable)
	photoCfg, ok := chattable.(tgbotapi.PhotoConfig)
	assert.True(t, ok)
	assert.NotNil(t, photoCfg)
	assert.EqualValues(t, tgbotapi.ModeMarkdown, photoCfg.ParseMode)
	assert.EqualValues(t, chatId, photoCfg.ChatID)
}

func TestUsecase_HandleFillingProfile_StageSexIncorrect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputText := "WrongSex"
	var chatId int64 = 1
	photoId := "photoId"
	user := &models.User{Id: "id", Stage: MaxProfileStage}

	usersRepo := mock.NewMockUsersRepository(ctrl)

	usecase := NewUsecase(
		usersRepo,
		nil,
		nil,
		zaptest.NewLogger(t).Sugar(),
	)

	chattable, err := usecase.HandleFillingProfile(context.Background(), inputText, chatId, photoId, user)
	assert.Nil(t, err)
	assert.NotNil(t, chattable)
	messageCfg, ok := chattable.(tgbotapi.MessageConfig)
	assert.True(t, ok)
	assert.NotNil(t, messageCfg)
	assert.EqualValues(t, chatId, messageCfg.ChatID)
}
