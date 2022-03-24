//go:generate mockgen -source users_repository.go -destination mock/users_repository.go -package mock
package internal

import (
	"context"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
)

type UsersRepository interface {
	Add(context.Context, *models.User) error
	GetByUserId(context.Context, string) (*models.User, error)
	UpdateByUserId(context.Context, *models.User) error
	DeleteByUserId(context.Context, string) error
	GetNextUser(context.Context, string, bool) (*models.User, error)
}
