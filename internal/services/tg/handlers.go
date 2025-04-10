package tg

import (
	"context"
	"fmt"
	"os"

	"github.com/un1uckyyy/email-in-tg/internal/models"

	"github.com/go-playground/validator/v10"
	tele "gopkg.in/telebot.v4"
)

var validate = validator.New()

func (t *TelegramService) help(c tele.Context) error {
	text, err := renderHTMLTemplate(helpTmpl, os.Getenv("TELEGRAM_SUPPORT"))
	if err != nil {
		msg := fmt.Sprintf("failed to render template: %v", err)
		logger.Error(msg)
		return c.Send(somethingWentWrong)
	}

	err = c.Send(text, tele.ModeHTML)
	if err != nil {
		msg := fmt.Sprintf("failed to send message: %v", err)
		logger.Error(msg)
		return c.Send(somethingWentWrong)
	}

	return nil
}

func (t *TelegramService) start(c tele.Context) error {
	ctx := context.Background()

	if c.Chat().Type != tele.ChatSuperGroup {
		return c.Send("Привет! Для начала добавь меня в свою группу\n" +
			"Важно, чтобы в ней были темы, тогда я смогу отправлять определенные письма в разные темы",
		)
	}

	// TODO add check that group already registered

	args := c.Args()
	if len(args) != 2 {
		return c.Send("Отправь команду /start в следующем формате:\n"+
			"/start 'email-address' 'password'\n"+
			"Как сгенерировать пароль смотри <a href=\"https://help.mail.ru/mail/mailer/password/\">здесь</a>",
			tele.ModeHTML,
		)
	}

	email, password := args[0], args[1]

	chat := c.Chat()
	group := models.Group{
		ID:    chat.ID,
		Type:  string(chat.Type),
		Title: chat.Title,
		Login: &models.EmailLogin{
			Email:    email,
			Password: password,
		},
		IsActive: true,
	}

	err := validate.Struct(group)
	if err != nil {
		msg := fmt.Sprintf("group validation error: %v", err)
		logger.Error(msg)
		return c.Send("Что-то пошло не так")
	}

	err = t.repo.CreateGroup(ctx, group)
	if err != nil {
		msg := fmt.Sprintf("group creation error: %v", err)
		logger.Error(msg)
		return c.Send("Что-то пошло не так")
	}

	t.pool.Register <- &group

	return c.Send("Отлично!\n" +
		"Теперь, чтобы добавить почту отправь /subscribe в нужную тему",
	)
}

func (t *TelegramService) stop(c tele.Context) error {
	return nil
}

func (t *TelegramService) login(c tele.Context) error {
	ctx := context.Background()

	args := c.Args()
	if len(args) != 2 {
		return c.Send("Отправь команду /login в следующем формате:\n"+
			"/login 'email-address' 'password'\n"+
			"Как сгенерировать пароль смотри <a href=\"https://help.mail.ru/mail/mailer/password/\">здесь</a>",
			tele.ModeHTML,
		)
	}

	groupID := c.Chat().ID
	email, password := args[0], args[1]
	credentials := models.EmailLogin{
		Email:    email,
		Password: password,
	}

	err := validate.Struct(credentials)
	if err != nil {
		msg := fmt.Sprintf("email validation error: %v", err)
		logger.Error(msg)
		return c.Send("Что-то пошло не так")
	}

	err = t.repo.SetEmailLogin(ctx, groupID, credentials)
	if err != nil {
		msg := fmt.Sprintf("group imap credentials set error: %v", err)
		logger.Error(msg)
		return c.Send("Что-то пошло не так")
	}

	return c.Send("Login успешен!")
}

func (t *TelegramService) subscribe(c tele.Context) error {
	ctx := context.Background()

	args := c.Args()
	if len(args) != 1 {
		return c.Send("Отправь команду /subscribe в следующем формате:\n" +
			"/subscribe 'email-address'",
		)
	}

	email := args[0]
	subscription := models.Subscription{
		SenderEmail: &email,
		GroupID:     c.Chat().ID,
		ThreadID:    c.Message().ThreadID,
	}

	err := validate.Struct(subscription)
	if err != nil {
		msg := fmt.Sprintf("subscription validation error: %v", err)
		logger.Error(msg)
		return c.Send("Что-то пошло не так")
	}

	err = t.repo.CreateSubscription(ctx, subscription)
	if err != nil {
		msg := fmt.Sprintf("subscription set error: %v", err)
		logger.Error(msg)
		return c.Send("Что-то пошло не так")
	}

	return c.Send("Подписка на почту создана\n" +
		"Письма от: " + email + " будут приходить в эту тему",
	)
}

func (t *TelegramService) subscribeOthers(c tele.Context) error {
	ctx := context.Background()

	subscription := models.Subscription{
		SenderEmail:  nil,
		GroupID:      c.Chat().ID,
		ThreadID:     c.Message().ThreadID,
		OtherSenders: true,
	}

	err := validate.Struct(subscription)
	if err != nil {
		msg := fmt.Sprintf("subscription validation error: %v", err)
		logger.Error(msg)
		return c.Send("Что-то пошло не так")
	}

	err = t.repo.CreateSubscription(ctx, subscription)
	if err != nil {
		msg := fmt.Sprintf("subscription set error: %v", err)
		logger.Error(msg)
		return c.Send("Что-то пошло не так")
	}

	return c.Send("Подписка на почту создана\n" +
		"Письма от всех незарегистрированных отправителей будут приходить в эту тему",
	)
}

func (t *TelegramService) subscriptions(c tele.Context) error {
	return nil
}
