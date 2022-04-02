package main

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	"github.com/Eretic431/datingTelegramBot/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

var (
	stages   = make(map[int]string, 6)
	commands = make(map[string]struct{}, 3)
)

func (a *application) handleUpdates() {
	commands["start"] = struct{}{}
	commands["profile"] = struct{}{}
	commands["next"] = struct{}{}

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

	if len(msg.Photo) > 0 && msg.Photo[0].FileID != "" {
		photo := msg.Photo[0].FileID
		a.log.Infof("receive file with id = %s", photo)
	}

	if user != nil && user.Stage != usecase.ProfileStageNone {
		fileId := "-"
		if len(msg.Photo) > 0 && msg.Photo[0].FileID != "" {
			fileId = msg.Photo[0].FileID
		}

		if msg.IsCommand() {
			outputMsg = tgbotapi.NewMessage(msg.Chat.ID, "Пожалуйста дозаполните анкету.")
		} else {
			outputMsg, err = a.usecase.HandleFillingProfile(ctx, msg.Text, msg.Chat.ID, fileId, user)
			if err != nil {
				return
			}
		}
	} else {
		_, ok := commands[msg.Command()]
		if ok {
			started, err := a.usecase.IsStarted(ctx, msg)
			if err != nil {
				return
			}

			if started {
				var err error
				switch msg.Command() {
				case "start":
					outputMsg, err = a.usecase.HandleStart(ctx, msg, started)
				case "profile":
					outputMsg, err = a.usecase.HandleProfile(ctx, msg, user)
				case "next":
					outputMsg, err = a.usecase.HandleCommandNext(ctx, msg.Chat.ID, user)
				}
				if err != nil {
					return
				}
			} else {
				outputMsg, err = a.usecase.HandleStart(ctx, msg, started)
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

	var msg tgbotapi.Chattable
	var err error

	if strings.HasPrefix(cq.Data, "like") || strings.HasPrefix(cq.Data, "dislike") {
		splitedData := strings.Split(cq.Data, ";")
		userId := splitedData[1]

		likeValue := splitedData[0] == "like"

		oldLike, err := a.likes.Get(ctx, cq.From.UserName, userId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				err := a.likes.Add(ctx, &models.Like{
					FromId: cq.From.UserName,
					ToId:   userId,
					Value:  likeValue,
				})

				if err != nil {
					a.log.Errorf("could not insert like with error %e", err)
					return
				}
			} else {
				a.log.Errorf("could not get like with error %e", err)
				return
			}
		} else {
			oldLike.Value = likeValue
			err := a.likes.Update(ctx, oldLike)
			if err != nil {
				a.log.Errorf("could not update like with error %e", err)
				return
			}
		}

		if likeValue {

			reverseLike, err := a.likes.Get(ctx, userId, cq.From.UserName)
			if err != nil {
				if !errors.Is(err, models.ErrNoRecord) {
					a.log.Errorf("could not get reverse like with error %e", err)
					return
				}
			} else if reverseLike != nil && reverseLike.Value {
				user2, err := a.users.GetByUserId(ctx, reverseLike.FromId)
				if err != nil {
					a.log.Errorf("could not get user with error %e", err)
					return
				}

				ava1 := tgbotapi.FileID(user.Image)
				match2Message := tgbotapi.NewPhoto(user2.ChatId, ava1)
				match2Message.Caption = internal.CreateMatchCaption(user)
				match2Message.ParseMode = tgbotapi.ModeMarkdown

				if _, err := a.bot.Send(match2Message); err != nil {
					a.log.Errorf("could not send message with error %e", err)
					return
				}

				ava2 := tgbotapi.FileID(user2.Image)
				match1Message := tgbotapi.NewPhoto(user.ChatId, ava2)
				match1Message.Caption = internal.CreateMatchCaption(user2)
				match1Message.ParseMode = tgbotapi.ModeMarkdown

				if _, err := a.bot.Send(match1Message); err != nil {
					a.log.Errorf("could not send message with error %e", err)
					return
				}
			}
		}

		msg, err = a.usecase.HandleCommandNext(ctx, cq.Message.Chat.ID, user)
		if err != nil {
			return
		}

	} else {
		msg, err = a.usecase.HandleFillingProfile(ctx, cq.Data, cq.Message.Chat.ID, user.Image, user)
	}

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
