package tg

import (
	"context"
	"fmt"
	"time"

	"github.com/un1uckyyy/email-in-tg/internal/models"
	"github.com/un1uckyyy/email-in-tg/internal/repo"

	"github.com/un1uckyyy/email-in-tg/internal/services/pool"
	tele "gopkg.in/telebot.v4"
)

const (
	pollerTimeout = 10 * time.Second
)

type TelegramService struct {
	bot  *tele.Bot
	pool *pool.Pool
	repo *repo.Repo
}

func NewTelegramService(
	token string,
	p *pool.Pool,
	repo *repo.Repo,
) (*TelegramService, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: pollerTimeout},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &TelegramService{
		bot:  bot,
		pool: p,
		repo: repo,
	}, nil
}

func (t *TelegramService) Start(ctx context.Context) error {
	groups, err := t.repo.GetAllActiveGroups(ctx)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("got %v active groups", len(groups))
	logger.Debug(msg)

	for _, group := range groups {
		t.pool.Register <- group
	}
	msg = fmt.Sprintf("register all %v groups", len(groups))
	logger.Debug(msg)

	logger.Debug("starting tg bot...")
	go t.run()
	return nil
}

func (t *TelegramService) run() {
	t.registerButtons()
	t.bot.Start()
}

func (t *TelegramService) registerButtons() {
	t.bot.Handle("/start", t.start)
	t.bot.Handle("/login", t.login)
	t.bot.Handle("/subscribe", t.subscribe)
	t.bot.Handle("/fetch", t.fetch)
}

func (t *TelegramService) Stop() {
	t.bot.Stop()
}

func (t *TelegramService) Send(groupID int64, threadID int, email *models.Email) error {
	logger.Debug(fmt.Sprintf("readers len: %v", len(email.Files)))

	group := &tele.User{ID: groupID}
	_, err := t.bot.Send(group, cleanTelegramHTML(email.Text), &tele.SendOptions{ThreadID: threadID}, tele.ModeHTML)
	if err != nil {
		logger.Error(err.Error())
	}

	media := make(tele.Album, 0, len(email.Files))
	for _, file := range email.Files {
		media = append(media,
			&tele.Document{
				File:     tele.FromReader(file.Data),
				FileName: file.Filename,
			},
		)
	}

	_, err = t.bot.SendAlbum(group, media, &tele.SendOptions{ThreadID: threadID})
	if err != nil {
		return err
	}

	return nil
}
