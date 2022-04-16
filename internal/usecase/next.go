package usecase

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Usecase) HandleCommandNext(ctx context.Context, chatId int64, user *models.User) (tgbotapi.Chattable, error) {
	u.log.Info("handleCommandNext")

	nextUser, err := u.users.GetNextUser(ctx, user.Id, user.Sex)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return tgbotapi.NewMessage(chatId, "Все анкеты просмотрены. Попробуйте ещё раз немного позже."), nil
		}
		u.log.Errorf("could not get next user with error %e", err)
		return tgbotapi.MessageConfig{}, err
	}

	if nextUser != nil {
		if len(nextUser.Image) > 0 {
			photoCfg := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(nextUser.Image))
			photoCfg.Caption = internal.CreateProfileCaption(nextUser)
			photoCfg.ParseMode = tgbotapi.ModeMarkdown

			photoCfg.ReplyMarkup = internal.CreateLikeKeyboardMarkup(nextUser.Id)
			return photoCfg, nil
		} else {
			msgConfig := tgbotapi.NewMessage(chatId, internal.CreateProfileCaption(nextUser))
			msgConfig.ParseMode = tgbotapi.ModeMarkdown
			msgConfig.ReplyMarkup = internal.CreateLikeKeyboardMarkup(nextUser.Id)
			return msgConfig, nil
		}
	}

	return tgbotapi.MessageConfig{}, nil
}
