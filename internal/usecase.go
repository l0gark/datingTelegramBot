//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package internal

import (
	"context"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Usecase interface {
	HandleStart(context.Context, *tgbotapi.Message, bool) (tgbotapi.MessageConfig, error)
	IsStarted(context.Context, *tgbotapi.Message) (bool, error)
	HandleProfile(context.Context, *tgbotapi.Message, *models.User) (tgbotapi.MessageConfig, error)
	HandleFillingProfile(context.Context, string, int64, string, *models.User) (tgbotapi.Chattable, error)
	HandleCommandNext(context.Context, int64, *models.User) (tgbotapi.Chattable, error)
}
