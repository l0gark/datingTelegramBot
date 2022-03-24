package usecase

import (
	"github.com/Eretic431/datingTelegramBot/internal"
	"github.com/Eretic431/datingTelegramBot/internal/data/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Usecase struct {
	users  internal.UserRepository
	likes  *postgres.LikeRepository
	bot    *tgbotapi.BotAPI
	log    *zap.SugaredLogger
	stages map[int]string
}

var _ internal.Usecase = &Usecase{}

func NewUsecase(
	users internal.UserRepository,
	likes *postgres.LikeRepository,
	bot *tgbotapi.BotAPI,
	log *zap.SugaredLogger) *Usecase {
	stages := make(map[int]string, 6)

	stages[0] = "Как Вас зовут?"
	stages[1] = "Сколько Вам лет?"
	stages[2] = "Из какого Вы города?"
	stages[3] = "Введите краткое описание своего профиля."
	stages[4] = "Пришлите фотографию, которая будет показываться другим пользователям в ленте."
	stages[5] = "Какого Вы пола? М/Ж"

	return &Usecase{
		users:  users,
		likes:  likes,
		bot:    bot,
		log:    log,
		stages: stages,
	}
}

const (
	MaxProfileStage  = 5
	ProfileStageNone = -1
)
