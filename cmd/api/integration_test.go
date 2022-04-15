package main

import (
	"context"
	"fmt"
	"github.com/Eretic431/datingTelegramBot/internal/data/models"
	"github.com/Eretic431/datingTelegramBot/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Scenario1_1(t *testing.T) {
	app := newTestApp()

	msg := &tgbotapi.Message{
		From:     &tgbotapi.User{UserName: "test"},
		Text:     "/start",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: 6}},
		Chat:     &tgbotapi.Chat{ID: 1},
	}
	ctx := context.Background()
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := fmt.Sprintf("Привет! Я, %s, помогаю людям познакомиться\n\n"+
		"*Список доступных команд:* \n"+
		"- /start - начало работы\n"+
		"- /profile - заполнить анкету\n"+
		"- /next - показать следующего пользователя",
		app.bot.Self.UserName,
	)

	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario1_2(t *testing.T) {
	app := newTestApp()

	msg := &tgbotapi.Message{
		From:     &tgbotapi.User{UserName: "test"},
		Text:     "/profile",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: 8}},
		Chat:     &tgbotapi.Chat{ID: 1},
	}
	ctx := context.Background()
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := fmt.Sprintf("Привет! Я, %s, помогаю людям познакомиться\n\n"+
		"*Список доступных команд:* \n"+
		"- /start - начало работы\n"+
		"- /profile - заполнить анкету\n"+
		"- /next - показать следующего пользователя",
		app.bot.Self.UserName,
	)

	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario2(t *testing.T) {
	app := newTestApp()

	msg := &tgbotapi.Message{
		From:     &tgbotapi.User{UserName: "test"},
		Text:     "/start",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: 6}},
		Chat:     &tgbotapi.Chat{ID: 1},
	}
	ctx := context.Background()
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := fmt.Sprintf("Привет! Я, %s, помогаю людям познакомиться\n\n"+
		"*Список доступных команд:* \n"+
		"- /start - начало работы\n"+
		"- /profile - заполнить анкету\n"+
		"- /next - показать следующего пользователя",
		app.bot.Self.UserName,
	)

	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario3(t *testing.T) {
	app := newTestApp()

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/start"}
	ctx := context.Background()
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := fmt.Sprintf("Привет! Я, %s, помогаю людям познакомиться\n\n"+
		"*Список доступных команд:* \n"+
		"- /start - начало работы\n"+
		"- /profile - заполнить анкету\n"+
		"- /next - показать следующего пользователя",
		app.bot.Self.UserName,
	)

	assert.Equal(t, expected, resp.Text)

	msg = &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/start"}
	chattable, _ = app.handleMessage(ctx, msg)
	resp = chattable[0].(tgbotapi.MessageConfig)

	expected = "Вы уже зарегистрированы в системе"

	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario4(t *testing.T) {
	app := newTestApp()

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/start"}
	ctx := context.Background()
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := fmt.Sprintf("Привет! Я, %s, помогаю людям познакомиться\n\n"+
		"*Список доступных команд:* \n"+
		"- /start - начало работы\n"+
		"- /profile - заполнить анкету\n"+
		"- /next - показать следующего пользователя",
		app.bot.Self.UserName,
	)
	assert.Equal(t, expected, resp.Text)

	msg = &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/profile"}
	chattable, _ = app.handleMessage(ctx, msg)
	resp = chattable[0].(tgbotapi.MessageConfig)

	expected = app.usecase.(*usecase.Usecase).Stages[0]
	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario5(t *testing.T) {
	app := newTestApp()

	ctx := context.Background()
	user1 := &models.User{
		Name:        "Masha",
		Sex:         false,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user1)

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/profile"}
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := app.usecase.(*usecase.Usecase).Stages[0]
	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario6(t *testing.T) {
	app := newTestApp()

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/start"}
	ctx := context.Background()
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := fmt.Sprintf("Привет! Я, %s, помогаю людям познакомиться\n\n"+
		"*Список доступных команд:* \n"+
		"- /start - начало работы\n"+
		"- /profile - заполнить анкету\n"+
		"- /next - показать следующего пользователя",
		app.bot.Self.UserName,
	)
	assert.Equal(t, expected, resp.Text)

	msg = &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/profile"}
	chattable, _ = app.handleMessage(ctx, msg)
	resp = chattable[0].(tgbotapi.MessageConfig)

	expected = app.usecase.(*usecase.Usecase).Stages[0]
	assert.Equal(t, expected, resp.Text)

	msg = &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Document: &tgbotapi.Document{}}
	chattable, _ = app.handleMessage(ctx, msg)

	expected = "Данные введены некорректно, попробуйте снова."
	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario7(t *testing.T) {
	app := newTestApp()

	ctx := context.Background()
	user1 := &models.User{
		Name:        "Masha",
		Sex:         false,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user1)

	user2 := &models.User{
		Name:        "Arkasha",
		Sex:         true,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user2)

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/next"}
	chattable, _ := app.handleMessage(ctx, msg)
	_, ok := chattable[0].(*tgbotapi.PhotoConfig)

	assert.True(t, ok)
}

func Test_Scenario8(t *testing.T) {
	app := newTestApp()

	ctx := context.Background()
	user1 := &models.User{
		Name:        "Masha",
		Sex:         false,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user1)

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/next"}
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := "Все анкеты просмотрены. Попробуйте ещё раз немного позже."

	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario9(t *testing.T) {
	app := newTestApp()

	ctx := context.Background()
	user1 := &models.User{
		Name:        "Masha",
		Sex:         false,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user1)

	user2 := &models.User{
		Name:        "Arkasha",
		Sex:         true,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user2)

	cq := &tgbotapi.CallbackQuery{From: &tgbotapi.User{UserName: "Masha"}, Data: "like;Arkasha"}
	_, _ = app.handleCallbackQuery(ctx, cq)

	cq = &tgbotapi.CallbackQuery{From: &tgbotapi.User{UserName: "Arkasha"}, Data: "like;Masha"}
	chattable, _ := app.handleCallbackQuery(ctx, cq)
	_, ok := chattable[0].(*tgbotapi.PhotoConfig)
	assert.True(t, ok)
	_, ok = chattable[1].(*tgbotapi.PhotoConfig)
	assert.True(t, ok)
}

func Test_Scenario10(t *testing.T) {
	app := newTestApp()

	ctx := context.Background()
	user1 := &models.User{
		Name:        "Masha",
		Sex:         false,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user1)

	user2 := &models.User{
		Name:        "Arkasha",
		Sex:         true,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user2)

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/next"}
	chattable, _ := app.handleMessage(ctx, msg)
	_, ok := chattable[0].(tgbotapi.MessageConfig)

	assert.True(t, ok)
}

func Test_Scenario11(t *testing.T) {
	app := newTestApp()

	ctx := context.Background()
	user1 := &models.User{
		Name:        "Masha",
		Sex:         false,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user1)

	user2 := &models.User{
		Name:        "Arkasha",
		Sex:         true,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user2)

	cq := &tgbotapi.CallbackQuery{From: &tgbotapi.User{UserName: "Masha"}, Data: "like;Arkasha"}
	_, _ = app.handleCallbackQuery(ctx, cq)

	_, err := app.likes.Get(ctx, "Masha", "Arkasha")
	assert.Nil(t, err)
}

func Test_Scenario12(t *testing.T) {
	app := newTestApp()

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/start"}
	ctx := context.Background()
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := fmt.Sprintf("Привет! Я, %s, помогаю людям познакомиться\n\n"+
		"*Список доступных команд:* \n"+
		"- /start - начало работы\n"+
		"- /profile - заполнить анкету\n"+
		"- /next - показать следующего пользователя",
		app.bot.Self.UserName,
	)

	assert.Equal(t, expected, resp.Text)

	msg = &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/next"}
	chattable, _ = app.handleMessage(ctx, msg)
	resp = chattable[0].(tgbotapi.MessageConfig)

	expected = "Все анкеты просмотрены. Попробуйте ещё раз немного позже."

	assert.Equal(t, expected, resp.Text)

	msg = &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/profile"}
	ctx = context.Background()
	chattable, _ = app.handleMessage(ctx, msg)
	resp = chattable[0].(tgbotapi.MessageConfig)

	expected = app.usecase.(*usecase.Usecase).Stages[0]

	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario13(t *testing.T) {
	app := newTestApp()

	ctx := context.Background()
	user1 := &models.User{
		Name:        "Masha",
		Sex:         false,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       4,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user1)

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/profile"}
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := "Данные введены некорректно, попробуйте снова."

	assert.Equal(t, expected, resp.Text)
}

func Test_Scenario14(t *testing.T) {
	app := newTestApp()

	ctx := context.Background()
	user1 := &models.User{
		Name:        "Masha",
		Sex:         false,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user1)

	user2 := &models.User{
		Name:        "Arkasha",
		Sex:         true,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user2)

	user3 := &models.User{
		Name:        "Vitya",
		Sex:         true,
		Age:         20,
		Description: "haha",
		City:        "test",
		Image:       "hardcoded",
		Started:     true,
		Stage:       -1,
		ChatId:      123,
	}
	_ = app.users.Add(ctx, user3)

	cq := &tgbotapi.CallbackQuery{From: &tgbotapi.User{UserName: "Vitya"}, Data: "like;Masha"}
	_, _ = app.handleCallbackQuery(ctx, cq)

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "Masha"}, Text: "/next"}
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(*tgbotapi.PhotoConfig)

	expected := fmt.Sprintf("*Имя:* %s\n"+
		"*Возраст:* %d\n"+
		"*Город:* %s\n"+
		"*Описание:* %s\n"+
		"*Пол:* %s", "Vitya", 20, "test", "haha", "Мужчина")

	assert.Equal(t, expected, resp.Caption)
}

func Test_Scenario15(t *testing.T) {
	app := newTestApp()

	msg := &tgbotapi.Message{From: &tgbotapi.User{UserName: "test"}, Text: "/asdkajsd"}
	ctx := context.Background()
	chattable, _ := app.handleMessage(ctx, msg)
	resp := chattable[0].(tgbotapi.MessageConfig)

	expected := "Такой команды не существует.\n\n" +
		"*Список доступных команд:* \n" +
		"- /start - начало работы\n" +
		"- /profile - заполнить анкету\n" +
		"- /next - показать следующего пользователя"

	assert.Equal(t, expected, resp.Text)
}
