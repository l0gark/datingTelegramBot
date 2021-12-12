package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

const (
	MaxProfileStage  = 5
	ProfileStageNone = -1
)

var stages map[int]string = make(map[int]string, 6)

func (a *application) handleUpdates() {
	commands := make(map[string]struct{}, 3)
	commands["start"] = struct{}{}
	commands["profile"] = struct{}{}
	commands["next"] = struct{}{}
	commands["skip"] = struct{}{}

	stages[0] = "Как Вас зовут?"
	stages[1] = "Сколько Вам лет?"
	stages[2] = "Из какого Вы города?"
	stages[3] = "Введите краткое описание своего профиля."
	stages[4] = "Пришлите фотографию, которая будет показываться другим пользователям в ленте."
	stages[5] = "Какого Вы пола? М/Ж"

	for update := range a.updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		ctx := context.Background()

		user, err := a.users.GetByUserId(ctx, update.Message.From.UserName)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				continue
			}
		}

		var outputMsg tgbotapi.MessageConfig

		a.log.Infof("update with username = %s", update.Message.From.UserName)

		if user != nil && user.Stage != ProfileStageNone {
			a.log.Infof("update stage № = %d", user.Stage)

			queryData := ""
			if update.CallbackQuery != nil {
				queryData = update.CallbackQuery.Data
			}

			outputMsg, err = a.handleFillingProfile(ctx, update.Message, &queryData, user)
			if err != nil {
				return
			}
		} else {
			a.log.Infof("Handle message or command")

			_, ok := commands[update.Message.Command()]

			if ok {
				started, err := a.isStarted(ctx, update.Message)
				if err != nil {
					continue
				}

				if started {
					var err error

					switch update.Message.Command() {
					case "start":
						outputMsg, err = a.handleCommandStart(ctx, update.Message, started)
					case "profile":
						outputMsg, err = a.handleCommandProfile(ctx, update.Message, user)
					case "next":
						outputMsg, err = a.handleCommandNext(update.Message)
					}
					if err != nil {
						continue
					}
				} else {
					outputMsg, err = a.handleCommandStart(ctx, update.Message, started)
					if err != nil {
						continue
					}
				}
			} else {
				outputMsg = a.handleUndefinedMessage(update.Message)
			}
		}

		if update.CallbackQuery != nil {
			a.log.Infof("update.CallbackQuery != nil, data = %s", update.CallbackQuery.Data)
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := a.bot.Request(callback); err != nil {
				a.log.Errorf("could not request callback with error = %e", err)
				continue
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := a.bot.Send(msg); err != nil {
				panic(err)
			}
		}

		if _, err := a.bot.Send(outputMsg); err != nil {
			a.log.Warnf("could not send message with error %e", err)
			return
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

func (a *application) handleCommandStart(ctx context.Context, inputMsg *tgbotapi.Message, started bool) (tgbotapi.MessageConfig, error) {
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

func (a *application) handleCommandProfile(ctx context.Context, inputMsg *tgbotapi.Message, user *models.User) (tgbotapi.MessageConfig, error) {
	user.Stage = 0
	err := a.users.UpdateByUserId(ctx, user)
	if err != nil {
		a.log.Errorf("could not update user with error %e", err)
		return tgbotapi.MessageConfig{}, err
	}

	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, stages[user.Stage])
	outputMsg.ReplyMarkup = createSkipKeyboardMarkup(user.Name)

	return outputMsg, err
}

func (a *application) handleFillingProfile(ctx context.Context, inputMsg *tgbotapi.Message, callbackData *string, user *models.User) (tgbotapi.MessageConfig, error) {
	a.log.Info("handleFillingProfile")
	var text string
	skipData := "google"

	currentData := ""
	if len(inputMsg.Text) > 0 {
		currentData = inputMsg.Text
	} else if callbackData != nil && len(*callbackData) > 0 {
		currentData = *callbackData
	}

	correct := true

	// name, age, city, ...
	switch user.Stage {
	case 0:
		name := currentData
		if len(name) > 0 {
			user.Name = name
			skipData = strconv.Itoa(user.Age)
		} else {
			correct = false
			skipData = user.Name
		}
	case 1:
		age, err := strconv.Atoi(currentData)
		if err == nil {
			user.Age = age
			skipData = user.City
		} else {
			correct = false
			skipData = strconv.Itoa(user.Age)
		}
	case 2:
		city := currentData
		if len(city) > 0 {
			user.City = city
		} else {
			correct = false
			skipData = city
		}
	}

	if correct {
		if user.Stage < MaxProfileStage {
			user.Stage += 1
		} else {
			user.Stage = ProfileStageNone
		}

		err := a.users.UpdateByUserId(ctx, user)
		if err != nil {
			a.log.Errorf("could not update user with error %e", err)
			return tgbotapi.MessageConfig{}, err
		}

		text = stages[user.Stage]
	} else {
		text = "Данные введены некорректны, попробуйте снова."
	}

	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, text)

	outputMsg.ReplyMarkup = createSkipKeyboardMarkup(skipData)

	return outputMsg, nil
}

func (a *application) handleCommandNext(inputMsg *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	a.log.Info("handleCommandNext")
	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "Next")

	return outputMsg, nil
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
