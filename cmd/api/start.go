package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *application) handleCommandStart(ctx context.Context, inputMsg *tgbotapi.Message, started bool) (tgbotapi.MessageConfig, error) {
	var text string

	if started {
		text = "Вы уже зарегистрированы в системе"
	} else {
		text = fmt.Sprintf("Привет! Я, %s, помогаю людям познакомиться\n\n"+
			"*Список доступных команд:* \n"+
			"- /start - начало работы\n"+
			"- /profile - заполнить анкету\n"+
			"- /next - показать следующего пользователя",
			a.bot.Self.UserName,
		)

		user := &models.User{
			Id:          inputMsg.From.UserName,
			Name:        "",
			Sex:         false,
			Age:         0,
			Description: "",
			City:        "",
			Image:       "",
			Started:     true,
			Stage:       ProfileStageNone,
			ChatId:      inputMsg.Chat.ID,
		}

		err := a.users.UpdateByUserId(ctx, user)
		if err != nil {
			a.log.Errorf("could not update user with error %e", err)
			return tgbotapi.MessageConfig{}, err
		}
	}

	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, text)
	outputMsg.ParseMode = tgbotapi.ModeMarkdown

	return outputMsg, nil
}

func (a *application) isStarted(ctx context.Context, inputMsg *tgbotapi.Message) (bool, error) {
	user, err := a.users.GetByUserId(ctx, inputMsg.From.UserName)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			user := &models.User{
				Id:          inputMsg.From.UserName,
				Name:        "",
				Sex:         false,
				Age:         0,
				Description: "",
				City:        "",
				Image:       "",
				Started:     false,
				Stage:       ProfileStageNone,
				ChatId:      inputMsg.Chat.ID,
			}

			err := a.users.Add(ctx, user)
			if err != nil {
				a.log.Errorf("could not insert user with error %e", err)
				return false, err
			}
			return false, nil
		}
		return false, err
	}

	return user.Started, nil
}
