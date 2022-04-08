package usecase

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
)

func (u *Usecase) GetUserByIdOrNil(ctx context.Context, userId string) (*models.User, error) {
	user, err := u.users.GetByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return nil, nil
		}

		u.log.Errorf("could not get user with error %e", err)
		return nil, err
	}
	return user, nil
}
