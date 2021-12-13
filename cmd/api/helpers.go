package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
	likeData := tgbotapi.NewInlineKeyboardButtonData("❤", "like;" + toId)
	dislikeData := tgbotapi.NewInlineKeyboardButtonData("➡", "dislike;" + toId)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(likeData, dislikeData),
	)
}
