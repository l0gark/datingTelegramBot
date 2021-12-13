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
	var outputMsg tgbotapi.MessageConfig
	var err error
	if user != nil && user.Stage != ProfileStageNone {
		fileId := "-"
		if len(msg.Photo) > 0 && msg.Photo[0].FileID != "" {
			fileId = msg.Photo[0].FileID
		}
		outputMsg, err = a.handleFillingProfile(ctx, msg.Text, msg.Chat.ID, fileId, user)
		if err != nil {
			return
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
					outputMsg, err = a.handleCommandNext(msg)
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
	if len(user.Name) > 0 {
		outputMsg.ReplyMarkup = createSkipKeyboardMarkup(user.Name)
	}

	return outputMsg, err
}

func (a *application) handleFillingProfile(
	ctx context.Context,
	inputText string,
	chatId int64,
	photoId string,
	user *models.User) (tgbotapi.MessageConfig, error) {
	var text string
	skipData := ""

	currentData := ""
	if len(inputText) > 0 {
		currentData = inputText
	}

	correct := true

	// name, age, city, description, image, sex
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
			skipData = user.Description
		} else {
			correct = false
			skipData = user.City
		}
	case 3:
		description := currentData
		if len(description) > 0 {
			user.Description = description
			if len(user.Image) == 0 {
				skipData = "emptyImage"
			} else {
				skipData = "NoImageData"
			}
		} else {
			correct = false
			skipData = user.Description
		}
	case 4:
		if photoId == "NoImageData" {
			if user.Sex {
				skipData = "М"
			} else {
				skipData = "Ж"
			}
		} else {
			if photoId != "-" {
				user.Image = photoId
				if user.Sex {
					skipData = "М"
				} else {
					skipData = "Ж"
				}
			} else {
				correct = false
				skipData = user.Image
			}
		}
	case 5:
		sex := currentData
		if sex == "М" || sex == "Ж" {
			if sex == "М" {
				user.Sex = true
			} else {
				user.Sex = false
			}
		} else {
			correct = false
		}
	case 6:
		skipData = ""

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

		if user.Stage == ProfileStageNone {
			text = "Профиль заполнен успешно.\nПопробуйте ввести команду /next"
		} else {
			text = stages[user.Stage]
		}
	} else {
		text = "Данные введены некорректно, попробуйте снова."
	}

	outputMsg := tgbotapi.NewMessage(chatId, text)

	if len(skipData) > 0 && skipData != "emptyImage" {
		outputMsg.ReplyMarkup = createSkipKeyboardMarkup(skipData)
	}

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
