package tg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"
	"github.com/un1uckyyy/email-in-tg/internal/domain/repo"

	"github.com/un1uckyyy/email-in-tg/internal/app/pool"
	tele "gopkg.in/telebot.v4"
)

const (
	pollerTimeout = 10 * time.Second
)

type TelegramService interface {
	Start(ctx context.Context) error
	Stop()
}

type telegramService struct {
	bot       *tele.Bot
	pool      pool.Pool
	groupRepo repo.GroupRepository
	subRepo   repo.SubscriptionRepository
}

func NewTelegramService(
	token string,
	p pool.Pool,
	groupRepo repo.GroupRepository,
	subRepo repo.SubscriptionRepository,
) (TelegramService, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: pollerTimeout},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &telegramService{
		bot:       bot,
		pool:      p,
		groupRepo: groupRepo,
		subRepo:   subRepo,
	}, nil
}

func (t *telegramService) Start(ctx context.Context) error {
	groups, err := t.groupRepo.GetAllActiveGroups(ctx)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("got %v active groups", len(groups))
	logger.Debug(msg)

	for _, group := range groups {
		err := t.pool.Add(ctx, group)
		if err != nil {
			return fmt.Errorf("failed to add group %v to pool: %v", group, err)
		}
	}
	msg = fmt.Sprintf("register all %v groups", len(groups))
	logger.Debug(msg)

	t.registerButtons()

	go func() { t.bot.Start() }()

	logger.Debug("starting tg bot...")
	go t.run(ctx)
	return nil
}

func (t *telegramService) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-t.pool.Updates():
			sub, err := t.subRepo.FindSubscription(ctx, update.GroupID, update.Email.MailFrom)
			if errors.Is(err, repo.ErrSubscriptionNotFound) {
				msg := fmt.Sprintf("subscription for email %v not found", update.Email.MailFrom)
				logger.Error(msg)
				break
			}
			if err != nil {
				msg := fmt.Sprintf("failed to find subscription: %v", err)
				logger.Error(msg)
				break
			}

			msg := fmt.Sprintf("found subscription: %+v", sub)
			logger.Debug(msg)

			err = t.send(ctx, sub.GroupID, sub.ThreadID, update.Email)
			if err != nil {
				msg := fmt.Sprintf("failed to send email: %v", err)
				logger.Error(msg)
				break
			}
		}
	}
}

func (t *telegramService) registerButtons() {
	t.bot.Handle("/help", t.help)
	t.bot.Handle("/start", t.start)
	t.bot.Handle("/stop", t.stop)
	t.bot.Handle("/login", t.login)
	t.bot.Handle("/subscribe", t.subscribe)
	t.bot.Handle("/subscribe_others", t.subscribeOthers)
	t.bot.Handle("/subscriptions", t.subscriptions)
}

func (t *telegramService) Stop() {
	t.bot.Stop()
}

func (t *telegramService) send(_ context.Context, groupID int64, threadID int, email *models.Email) error {
	logger.Debug(fmt.Sprintf("readers len: %v", len(email.Files)))

	group := &tele.User{ID: groupID}
	email.Text = cleanTelegramHTML(email.Text)
	text, err := renderHTMLTemplate(emailTmpl, email)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	_, err = t.bot.Send(
		group,
		text,
		&tele.SendOptions{ThreadID: threadID},
		tele.ModeHTML,
	)
	if err != nil {
		return fmt.Errorf("failed to send email text: %w", err)
	}

	albums := splitFilesToAlbums(email.Files)
	for _, album := range albums {
		_, err = t.bot.SendAlbum(group, album, &tele.SendOptions{ThreadID: threadID})
		if err != nil {
			return fmt.Errorf("failed to send album: %w", err)
		}
	}

	return nil
}
