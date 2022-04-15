//go:generate mockgen -source likes_repository.go -destination mock/likes_repository.go -package mock
package internal

import (
	"context"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
)

type LikesRepository interface {
	Add(context.Context, *models.Like) error
	Get(context.Context, string, string) (*models.Like, error)
	Update(context.Context, *models.Like) error
	Delete(context.Context, int64) error
	DeleteAll(ctx context.Context) error
}
