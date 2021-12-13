package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *application) handleCommandNext(ctx context.Context, chatId int64, user *models.User) (tgbotapi.Chattable, error) {
	a.log.Info("handleCommandNext")

	nextUser, err := a.users.GetNextUser(ctx, user.Id, user.Sex)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return tgbotapi.NewMessage(chatId, "А уже всё!"), nil
		}
		a.log.Errorf("could not get next user with error %e", err)
		return tgbotapi.MessageConfig{}, err
	}

	if nextUser != nil {
		photoCfg := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(nextUser.Image))
		sex := ""
		if nextUser.Sex {
			sex = "Мужчина"
		} else {
			sex = "Женщина"
		}
		caption := fmt.Sprintf("*Имя:* %s\n"+
			"*Возраст:* %d\n"+
			"*Город:* %s\n"+
			"*Описание:* %s\n"+
			"*Пол:* %s", nextUser.Name, nextUser.Age, nextUser.City, nextUser.Description, sex)
		photoCfg.Caption = caption
		photoCfg.ParseMode = tgbotapi.ModeMarkdown

		photoCfg.ReplyMarkup = createLikeKeyboardMarkup(nextUser.Id)
		return photoCfg, nil
	}

	return tgbotapi.MessageConfig{}, nil
}
