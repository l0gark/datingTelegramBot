package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Usecase) HandleStart(
	ctx context.Context,
	inputMsg *tgbotapi.Message,
	started bool) (tgbotapi.MessageConfig, error) {
	var text string

	if started {
		text = "Вы уже зарегистрированы в системе"
	} else {
		text = fmt.Sprintf("Привет! Я, %s, помогаю людям познакомиться\n\n"+
			"*Список доступных команд:* \n"+
			"- /start - начало работы\n"+
			"- /profile - заполнить анкету\n"+
			"- /next - показать следующего пользователя",
			u.bot.Self.UserName,
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

		err := u.users.UpdateByUserId(ctx, user)
		if err != nil {
			u.log.Errorf("could not update user with error %e", err)
			return tgbotapi.MessageConfig{}, err
		}
	}

	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, text)
	outputMsg.ParseMode = tgbotapi.ModeMarkdown

	return outputMsg, nil
}

func (u *Usecase) IsStarted(ctx context.Context, inputMsg *tgbotapi.Message) (bool, error) {
	user, err := u.users.GetByUserId(ctx, inputMsg.From.UserName)
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

			err := u.users.Add(ctx, user)
			if err != nil {
				u.log.Errorf("could not insert user with error %e", err)
				return false, err
			}
			return false, nil
		}
		return false, err
	}

	return user.Started, nil
}
