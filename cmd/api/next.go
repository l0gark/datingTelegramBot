package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (a *application) handleCommandNext(inputMsg *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	a.log.Info("handleCommandNext")
	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "Next")

	return outputMsg, nil
}
