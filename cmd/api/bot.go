package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (a *application) handleUpdates() {
	commands := make(map[string]struct{}, 3)
	commands["/start"] = struct{}{}
	commands["/profile"] = struct{}{}
	commands["/next"] = struct{}{}

	for update := range a.updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		_, ok := commands[update.Message.Text]

		var outputMsg *tgbotapi.MessageConfig

		if ok {
			ctx := context.Background()

			started, err := a.isStarted(ctx, update.Message)
			if err != nil {
				continue
			}

			if started {
				switch update.Message.Text {
				case "/start":
					outputMsg, _ = a.handleCommandStart(ctx, update.Message, true)
				case "/profile":
					outputMsg = a.handleCommandProfile(update.Message)
				case "/next":
					outputMsg = a.handleCommandNext(update.Message)
				}
			} else {
				outputMsg, err = a.handleCommandStart(ctx, update.Message, false)
				if err != nil {
					continue
				}
			}
		} else {
			outputMsg = a.handleUndefinedMessage(update.Message)
		}

		if outputMsg != nil {
			if _, err := a.bot.Send(outputMsg); err != nil {
				a.log.Warnf("could not send message with error %e", err)
				return
			}
		}
	}
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

func (a *application) handleCommandStart(ctx context.Context, inputMsg *tgbotapi.Message, started bool) (*tgbotapi.MessageConfig, error) {
	a.log.Infof("handleCommandStart, started = %t", started)

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
		}

		err := a.users.UpdateByUserId(ctx, user)
		if err != nil {
			a.log.Errorf("could not update user with error %e", err)
			return nil, err
		}
	}

	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, text)
	outputMsg.ParseMode = tgbotapi.ModeMarkdown

	return &outputMsg, nil
}

func (a *application) handleCommandProfile(inputMsg *tgbotapi.Message) *tgbotapi.MessageConfig {
	a.log.Info("handleCommandProfile")
	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "Hello")

	return &outputMsg
}

func (a *application) handleCommandNext(inputMsg *tgbotapi.Message) *tgbotapi.MessageConfig {
	a.log.Info("handleCommandNext")
	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "Hello")

	return &outputMsg
}

func (a *application) handleUndefinedMessage(inputMsg *tgbotapi.Message) *tgbotapi.MessageConfig {
	a.log.Info("handleUndefinedMessage")
	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "Такой команды не существует.\n\n"+
		"*Список доступных команд:* \n"+
		"- /start - начало работы\n"+
		"- /profile - заполнить анкету\n"+
		"- /next - показать следующего пользователя",
	)
	outputMsg.ParseMode = tgbotapi.ModeMarkdown

	return &outputMsg
}

func newTgBot(c *config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(c.TgBotToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = c.Production
	log.Printf("authorised on account %s\n", bot.Self.UserName)

	return bot, nil
}

func newTgBotUpdatesChan(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return updates
}
