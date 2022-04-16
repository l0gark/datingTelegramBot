package usecase

import (
	"context"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
)

func (u *Usecase) DeleteAll(ctx context.Context) error {
	if err := u.likes.DeleteAll(ctx); err != nil {
		u.log.Errorf("couldn't delete all likes with err = %e", err)
		return err
	}
	if err := u.users.DeleteAll(ctx); err != nil {
		u.log.Errorf("couldn't delete all users with err = %e", err)
		return err
	}
	return nil
}

func (u *Usecase) AddTestUser(ctx context.Context, sex bool) error {
	if err := u.users.Add(ctx, &models.User{
		Id:          "TestId",
		Name:        "TestName",
		Sex:         sex,
		Age:         0,
		Description: "TestDescription",
		City:        "TestCity",
		Image:       "",
		Started:     true,
		Stage:       -1,
		ChatId:      0,
	}); err != nil {
		u.log.Errorf("couldn't insert test user with err = %e", err)
		return err
	}
	return nil
}

func (u *Usecase) AddTestUserWithLike(ctx context.Context, sex bool, toId string) error {
	if err := u.AddTestUser(ctx, sex); err != nil {
		return err
	}
	if err := u.likes.Add(ctx, &models.Like{FromId: "TestId", ToId: toId, Value: true}); err != nil {
		u.log.Errorf("couldn't insert test user with err = %e", err)
		return err
	}
	return nil
}
