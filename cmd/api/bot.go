package main

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	MaxProfileStage  = 5
	ProfileStageNone = -1
)

var (
	stages   = make(map[int]string, 6)
	commands = make(map[string]struct{}, 3)
)

func (a *application) handleUpdates() {
	commands["start"] = struct{}{}
	commands["profile"] = struct{}{}
	commands["next"] = struct{}{}

	stages[0] = "Как Вас зовут?"
	stages[1] = "Сколько Вам лет?"
	stages[2] = "Из какого Вы города?"
	stages[3] = "Введите краткое описание своего профиля."
	stages[4] = "Пришлите фотографию, которая будет показываться другим пользователям в ленте."
	stages[5] = "Какого Вы пола? М/Ж"

	for update := range a.updates {
		ctx := context.Background()
		if update.Message != nil {
			user, err := a.users.GetByUserId(ctx, update.Message.From.UserName)
			if err != nil {
				if !errors.Is(err, models.ErrNoRecord) {
					continue
				}
			}
			a.handleMessages(ctx, update.Message, user)
		} else if update.CallbackQuery != nil {
			user, err := a.users.GetByUserId(ctx, update.CallbackQuery.From.UserName)
			if err != nil {
				if !errors.Is(err, models.ErrNoRecord) {
					continue
				}
			}
			a.handleCallbackQueries(ctx, update.CallbackQuery, user)
		}
	}
}

func (a *application) handleMessages(ctx context.Context, msg *tgbotapi.Message, user *models.User) {
	var outputMsg tgbotapi.Chattable
	var err error
	if user != nil && user.Stage != ProfileStageNone {
		fileId := "-"
		if len(msg.Photo) > 0 && msg.Photo[0].FileID != "" {
			fileId = msg.Photo[0].FileID
		}
		if msg.IsCommand() {
			outputMsg = tgbotapi.NewMessage(msg.Chat.ID, "Пожалуйста дозаполните анкету.")
		} else {
			outputMsg, err = a.handleFillingProfile(ctx, msg.Text, msg.Chat.ID, fileId, user)
			if err != nil {
				return
			}
		}
	} else {
		_, ok := commands[msg.Command()]
		if ok {
			started, err := a.isStarted(ctx, msg)
			if err != nil {
				return
			}

			if started {
				var err error
				switch msg.Command() {
				case "start":
					outputMsg, err = a.handleCommandStart(ctx, msg, started)
				case "profile":
					outputMsg, err = a.handleCommandProfile(ctx, msg, user)
				case "next":
					outputMsg, err = a.handleCommandNext(ctx, msg, user)
				}
				if err != nil {
					return
				}
			} else {
				outputMsg, err = a.handleCommandStart(ctx, msg, started)
				if err != nil {
					return
				}
			}
		} else {
			outputMsg = a.handleUndefinedMessage(msg)
		}
	}

	if _, err := a.bot.Send(outputMsg); err != nil {
		a.log.Warnf("could not send message with error %e", err)
		return
	}
}

func (a *application) handleCallbackQueries(ctx context.Context, cq *tgbotapi.CallbackQuery, user *models.User) {
	callback := tgbotapi.NewCallback(cq.ID, cq.Data)
	if _, err := a.bot.Request(callback); err != nil {
		a.log.Errorf("could not request callback with error %e", err)
		return
	}

	msg, err := a.handleFillingProfile(ctx, cq.Data, cq.Message.Chat.ID, user.Image, user)
	if err != nil {
		a.log.Errorf("could not handle profile filling with error %e", err)
		return
	}
	if _, err := a.bot.Send(msg); err != nil {
		a.log.Errorf("could not send message with error %e", err)
		return
	}
}

func (a *application) handleUndefinedMessage(inputMsg *tgbotapi.Message) tgbotapi.MessageConfig {
	a.log.Info("handleUndefinedMessage")
	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "Такой команды не существует.\n\n"+
		"*Список доступных команд:* \n"+
		"- /start - начало работы\n"+
		"- /profile - заполнить анкету\n"+
		"- /next - показать следующего пользователя",
	)
	outputMsg.ParseMode = tgbotapi.ModeMarkdown

	return outputMsg
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
