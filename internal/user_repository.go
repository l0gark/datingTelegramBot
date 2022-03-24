//go:generate mockgen -source user_repository.go -destination mock/user_repository.go -package mock
package internal

import (
	"context"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
)

type UserRepository interface {
	Add(ctx context.Context, user *models.User) error
	GetByUserId(ctx context.Context, userId string) (*models.User, error)
	UpdateByUserId(ctx context.Context, user *models.User) error
	DeleteByUserId(ctx context.Context, userId string) error
	GetNextUser(ctx context.Context, userId string, sex bool) (*models.User, error)
}
