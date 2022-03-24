//go:generate mockgen -source user_repository.go -destination mock/user_repository.go -package mock
package internal

import (
	"context"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
)

type UserRepository interface {
	Add(context.Context, *models.User) error
	GetByUserId(context.Context, string) (*models.User, error)
	UpdateByUserId(context.Context, *models.User) error
	DeleteByUserId(context.Context, string) error
	GetNextUser(context.Context, string, bool) (*models.User, error)
}
