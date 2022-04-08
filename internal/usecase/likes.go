package usecase

import (
	"context"
	"errors"
	"github.com/Eretic431/datingTelegramBot/internal"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Usecase) AddOrUpdateLike(ctx context.Context, likeValue bool, fromId, toId string) error {
	oldLike, err := u.likes.Get(ctx, fromId, toId)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			err := u.likes.Add(ctx, &models.Like{
				FromId: fromId,
				ToId:   toId,
				Value:  likeValue,
			})

			if err != nil {
				u.log.Errorf("could not insert like with error %e", err)
				return err
			}
		} else {
			u.log.Errorf("could not get like with error %e", err)
			return err
		}
	} else {
		oldLike.Value = likeValue
		err := u.likes.Update(ctx, oldLike)
		if err != nil {
			u.log.Errorf("could not update like with error %e", err)
			return err
		}
	}

	return nil
}

func (u *Usecase) HasLikeWithTrueValue(ctx context.Context, fromId, toId string) (bool, error) {
	reverseLike, err := u.likes.Get(ctx, fromId, toId)
	if err != nil {
		if !errors.Is(err, models.ErrNoRecord) {
			u.log.Errorf("could not get reverse like with error %e", err)
			return false, err
		}
		return false, nil
	}

	return reverseLike.Value, nil
}

func (u *Usecase) CreateMatchMessages(user1, user2 *models.User) (tgbotapi.Chattable, tgbotapi.Chattable) {
	ava1 := tgbotapi.FileID(user1.Image)
	match2Message := tgbotapi.NewPhoto(user2.ChatId, ava1)
	match2Message.Caption = internal.CreateMatchCaption(user1)
	match2Message.ParseMode = tgbotapi.ModeMarkdown

	ava2 := tgbotapi.FileID(user2.Image)
	match1Message := tgbotapi.NewPhoto(user1.ChatId, ava2)
	match1Message.Caption = internal.CreateMatchCaption(user2)
	match1Message.ParseMode = tgbotapi.ModeMarkdown

	return match1Message, match2Message
}
