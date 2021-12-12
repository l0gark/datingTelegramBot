package main

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (a *application) handleUpdates() {
	for update := range a.updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		ctx := context.Background()

		var outputMsg *tgbotapi.MessageConfig

		switch update.Message.Text {
		case "/start":
			outputMsg = a.handleCommandStart(update.Message)
		case "/profile":
			started, err := a.isStarted(ctx, update.Message)
			if err != nil {
				continue
			}

			if started {
				outputMsg = a.handleCommandProfile(update.Message)
			} else {
				outputMsg = a.handleCommandStart(update.Message)
			}
		case "/next":
			started, err := a.isStarted(ctx, update.Message)
			if err != nil {
				continue
			}

			if started {
				outputMsg = a.handleCommandNext(update.Message)
			} else {
				outputMsg = a.handleCommandStart(update.Message)
			}
		default:
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

func (a *application) handleCommandStart(inputMsg *tgbotapi.Message) *tgbotapi.MessageConfig {
	a.log.Info("handleCommandStart")
	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "Hello! I am dating bot")

	return &outputMsg
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
