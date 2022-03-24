package usecase

import (
	"context"
	"github.com/Eretic431/datingTelegramBot/internal"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func (u *Usecase) HandleProfile(
	ctx context.Context,
	inputMsg *tgbotapi.Message,
	user *models.User) (tgbotapi.MessageConfig, error) {
	user.Stage = 0
	err := u.users.UpdateByUserId(ctx, user)
	if err != nil {
		u.log.Errorf("could not update user with error %e", err)
		return tgbotapi.MessageConfig{}, err
	}

	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, u.stages[user.Stage])
	if len(user.Name) > 0 {
		outputMsg.ReplyMarkup = internal.CreateSkipKeyboardMarkup(user.Name)
	}

	return outputMsg, err
}
func (u *Usecase) HandleFillingProfile(
	ctx context.Context,
	inputText string,
	chatId int64,
	photoId string,
	user *models.User) (tgbotapi.Chattable, error) {
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

		err := u.users.UpdateByUserId(ctx, user)
		if err != nil {
			u.log.Errorf("could not update user with error %e", err)
			return tgbotapi.MessageConfig{}, err
		}

		if user.Stage != ProfileStageNone {
			text = u.stages[user.Stage]
		}
	} else {
		text = "Данные введены некорректно, попробуйте снова."
	}

	if user.Stage == ProfileStageNone {
		photCfg := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(user.Image))
		photCfg.Caption = internal.CreateMyProfileCaption(user)
		photCfg.ParseMode = tgbotapi.ModeMarkdown
		return photCfg, nil
	}

	outputMsg := tgbotapi.NewMessage(chatId, text)

	if len(skipData) > 0 && skipData != "emptyImage" {
		outputMsg.ReplyMarkup = internal.CreateSkipKeyboardMarkup(skipData)
	}

	return outputMsg, nil
}
