package main

import (
	"fmt"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createSkipKeyboardMarkup(data string) tgbotapi.InlineKeyboardMarkup {
	if len(data) == 0 {
		data = "-"
	}
	buttonData := tgbotapi.NewInlineKeyboardButtonData("Пропустить", data)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttonData),
	)
}

func createLikeKeyboardMarkup(toId string) tgbotapi.InlineKeyboardMarkup {
	likeData := tgbotapi.NewInlineKeyboardButtonData("❤", "like;"+toId)
	dislikeData := tgbotapi.NewInlineKeyboardButtonData("➡", "dislike;"+toId)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(likeData, dislikeData),
	)
}

func createMyProfileCaption(user *models.User) string {
	return createProfileCaption(user) + "\n\nПопробуйте ввести команду /next"
}

func createProfileCaption(user *models.User) string {
	sex := ""
	if user.Sex {
		sex = "Мужчина"
	} else {
		sex = "Женщина"
	}

	caption := fmt.Sprintf("*Имя:* %s\n"+
		"*Возраст:* %d\n"+
		"*Город:* %s\n"+
		"*Описание:* %s\n"+
		"*Пол:* %s", user.Name, user.Age, user.City, user.Description, sex)
	return caption
}

func createMatchCaption(user *models.User) string {
	return "Поздравляем! У Вас совпадание с @" + user.Id + "\nМожете связаться в личных сообщениях☺\n\n" + createProfileCaption(user)
}
