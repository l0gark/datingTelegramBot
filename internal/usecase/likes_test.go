package usecase

import (
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

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
