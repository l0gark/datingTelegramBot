package main

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *application) handleCommandNext(ctx context.Context, chatId int64, user *models.User) (tgbotapi.Chattable, error) {
	a.log.Info("handleCommandNext")

	nextUser, err := a.users.GetNextUser(ctx, user.Id, user.Sex)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return tgbotapi.NewMessage(chatId, "Все анкеты просмотрены. Попробуйте ещё раз немного позже."), nil
		}
		a.log.Errorf("could not get next user with error %e", err)
		return tgbotapi.MessageConfig{}, err
	}

	if nextUser != nil {
		photoCfg := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(nextUser.Image))
		photoCfg.Caption = createProfileCaption(nextUser)
		photoCfg.ParseMode = tgbotapi.ModeMarkdown

		photoCfg.ReplyMarkup = createLikeKeyboardMarkup(nextUser.Id)
		return photoCfg, nil
	}

	return tgbotapi.MessageConfig{}, nil
}
