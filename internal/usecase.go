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

	AddOrUpdateLike(ctx context.Context, likeValue bool, fromId, toId string) error
	HasLikeWithTrueValue(ctx context.Context, fromId, toId string) (bool, error)
	CreateMatchMessages(user1, user2 *models.User) (tgbotapi.Chattable, tgbotapi.Chattable)

	GetUserByIdOrNil(ctx context.Context, userId string) (*models.User, error)
}
