package tg

import (
	"context"
	"fmt"
	"strconv"

	"github.com/un1uckyyy/email-in-tg/internal/models"

	"github.com/go-playground/validator/v10"
	tele "gopkg.in/telebot.v4"
)

var validate = validator.New()

func (t *TelegramService) start(c tele.Context) error {
	ctx := context.Background()

	if c.Chat().Type != tele.ChatSuperGroup {
		return c.Send("Привет! Для начала добавь меня в свою группу\n" +
			"Важно, чтобы в ней были темы, тогда я смогу отправлять определенные письма в разные темы",
		)
	}

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

	return c.Send("Отлично!\n" +
		"Теперь, чтобы добавить почту отправь /subscribe в нужную тему",
	)
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
		SenderEmail: email,
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

func (t *TelegramService) fetch(c tele.Context) error {
	args := c.Args()
	if len(args) != 1 {
		return c.Send("Отправь команду /fetch в следующем формате:\n" +
			"/fetch 'number'",
		)
	}

	number := args[0]

	nu, err := strconv.ParseUint(number, 10, 32)
	if err != nil {
		logger.Error(err.Error())
	}

	ic, ok := t.pool.Clients[c.Chat().ID]
	if !ok {
		logger.Error("cant get imap client from pool")
		return nil
	}

	err = ic.Select("INBOX")
	if err != nil {
		logger.Error(err.Error())
		return c.Send("Что-то пошло не так")
	}

	email, err := ic.FetchOne(uint32(nu), true)
	if err != nil {
		logger.Error("fetch error: " + err.Error())
		return c.Send("Что-то пошло не так")
	}

	logger.Debug(fmt.Sprintf("groupId: %v, threadId: %v", c.Chat().ID, c.Message().ThreadID))
	err = t.Send(c.Chat().ID, c.Message().ThreadID, email)
	if err != nil {
		logger.Error(err.Error())
	}

	return c.Send("executed")
}
