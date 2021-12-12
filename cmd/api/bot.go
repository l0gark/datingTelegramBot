package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (a *application) handleUpdates() {
	for update := range a.updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		var outputMsg *tgbotapi.MessageConfig

		switch update.Message.Text {
		case "/start":
			outputMsg = a.handleCommandStart(update.Message)
		case "/profile":
			outputMsg = a.handleCommandProfile(update.Message)
		case "/next":
			outputMsg = a.handleCommandNext(update.Message)
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

func (a *application) handleCommandStart(inputMsg *tgbotapi.Message) *tgbotapi.MessageConfig {
	a.log.Info("handleCommandStart")
	outputMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "Hello")
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
